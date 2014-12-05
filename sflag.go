// Package sflag at time of writing, is the only known flags package variant that is 100% DRY, free of fugly pointer syntax and uses clean struct syntax.
// Implementation makes use of reflection and struct tags.
// Limitation: Presence of a boolean flag requires that there be no STANDALONE true or false parameters, use "--Foo=true" syntax instead of "--Foo true".
// This is because the underlying flags package will stop processing on seeing the first standalone true/false value
// (This is because it will considers the preceding bool flag (--Foo) set by its presence alone)
// Limitation: Commandline args must start with an upcase char, since the implementation uses reflection which fails on unexported fields
package sflag

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	visited map[string]bool
)

func noteVisited(_flag *flag.Flag) {
	visited[_flag.Name] = true
}

// Parse iterates through the members of the struct.
// Using type info obtained via reflection, and parsing the struct tag for usage and default value, it sets up golang's flag package to actually parse
// Normally, the rightmost pipe char in the tag is used to delineate between the Description (on the left) and Default value (on the right)
// (However, you can change delineator to the first char of the tag (after eliminating leading whitespace) if that char is not alphabetic)
// Fields with no tag or whitespace-only tags are ignored
// Non-nil pointer fields are ignored
// Nil pointer fields will be left nil if that flag is not set on commandline (and the tag is not parsed for a default value)
// Parameters not consumed by flags will be copied to the last field of type []string
func Parse(ss interface{}) {
	visited = make(map[string]bool)
	pointers := map[string]interface{}{}
	if reflect.TypeOf(ss).Kind() != reflect.Ptr {
		panic("sflag.Parse was not provided a pointer arg")
	}
	sstype := reflect.TypeOf(ss).Elem()
	ssvalue := reflect.ValueOf(ss).Elem()

	if sstype.Kind() != reflect.Struct {
		panic("sflag.Parse was not provided a pointer to a struct")
	}

	moreusage := ""

	var flags = *flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	var argsiface interface{}

	hasBoolArg := false

	for ii := 0; ii < sstype.NumField(); ii++ {
		pp := sstype.Field(ii)
		vv := ssvalue.Field(ii)
		if pp.Anonymous {
			continue
		}
		if pp.Name == "Usage" {
			continue
		}
		if pp.Type.String() == "[]string" {
			argsiface = vv.Addr().Interface()
			continue
		}
		tag := strings.TrimSpace((string)(pp.Tag))
		if tag == "" {
			continue
		}
		splitChar := tag[0:1]
		if strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", splitChar) {
			splitChar = "|"
		} else {
			tag = tag[1:]
		}
		parts := strings.Split(tag, splitChar)
		part0 := ""
		part1 := ""
		if len(parts) > 0 {
			part0 = strings.TrimSpace(parts[0])
		}
		if len(parts) > 1 {
			part1 = strings.TrimSpace(parts[1])
		}
		if len(parts) > 0 {
			moreusage += "\n\t--" + pp.Name + ": " + part1 + " <-- Default, " + pp.Type.String() + " # " + part0
		}
		if pp.Type.Kind() == reflect.Ptr {
			if vv.Elem().Kind() != reflect.Invalid {
				continue
			}
			switch pp.Type.String() {
			case "*string":
				tempstr := ""
				pointers[pp.Name] = &tempstr
				flags.StringVar(&tempstr, pp.Name, tempstr, "")
			case "*int":
				tempint := 0
				pointers[pp.Name] = &tempint
				flags.IntVar(&tempint, pp.Name, tempint, "")
			case "*bool":
				tempbool := false
				pointers[pp.Name] = &tempbool
				flags.BoolVar(&tempbool, pp.Name, tempbool, "")
			case "*int64":
				tempint64 := int64(0)
				pointers[pp.Name] = &tempint64
				flags.Int64Var(&tempint64, pp.Name, tempint64, "")
			case "*float64":
				tempfloat64 := 0.0
				pointers[pp.Name] = &tempfloat64
				flags.Float64Var(&tempfloat64, pp.Name, tempfloat64, "")
			default:
				continue
			}
		}

		if len(parts) == 1 {
			switch pp.Type.Kind() {
			case reflect.String:
				flags.StringVar(vv.Addr().Interface().(*string), pp.Name, vv.String(), " <--default, string # "+part0)
			case reflect.Int:
				flags.IntVar(vv.Addr().Interface().(*int), pp.Name, int(vv.Int()), " <--default, int # "+part0)
			case reflect.Bool:
				flags.BoolVar(vv.Addr().Interface().(*bool), pp.Name, bool(vv.Bool()), " <--default, bool # "+part0)
				hasBoolArg = true
			case reflect.Int64:
				flags.Int64Var(vv.Addr().Interface().(*int64), pp.Name, vv.Int(), " <--default, int64 # "+part0)
			case reflect.Float64:
				flags.Float64Var(vv.Addr().Interface().(*float64), pp.Name, vv.Float(), " <--default, float64 # "+part0)
			}
		}
		if len(parts) == 2 {
			switch pp.Type.Kind() {
			case reflect.String:
				vv.SetString(part1)
				flags.StringVar(vv.Addr().Interface().(*string), pp.Name, part1, " <--default, string # "+part0)
			case reflect.Int:
				inum, _ := strconv.ParseInt(part1, 10, 64)
				vv.SetInt(inum)
				flags.IntVar(vv.Addr().Interface().(*int), pp.Name, int(inum), " <--default, int # "+part0)
			case reflect.Bool:
				bnum, _ := strconv.ParseBool(part1)
				vv.SetBool(bnum)
				flags.BoolVar(vv.Addr().Interface().(*bool), pp.Name, bool(bnum), " <--default, bool # "+part0)
				hasBoolArg = true
			case reflect.Int64:
				jnum, _ := strconv.ParseInt(part1, 10, 64)
				vv.SetInt(jnum)
				flags.Int64Var(vv.Addr().Interface().(*int64), pp.Name, jnum, " <--default, int64 # "+part0)
			case reflect.Float64:
				fnum, _ := strconv.ParseFloat(part1, 64)
				vv.SetFloat(fnum)
				flags.Float64Var(vv.Addr().Interface().(*float64), pp.Name, fnum, " <--default, float64 # "+part0)
			}
		}
	}

	pp, _ := sstype.FieldByName("Usage")
	vv := ssvalue.FieldByName("Usage")
	vv.SetString("\n Usage of " + os.Args[0] + " # " + (string)(pp.Tag) + "\n ARGS:" + moreusage)
	if hasBoolArg {
		for _, arg := range os.Args[1:] {
			switch strings.ToLower(arg) {
			case "true", "false":
				panic("Golang flag package requires \"--Foo=bar\" instead of \"--Foo bar\" syntax for bool args")
			}
		}
	}
	flags.Parse(os.Args[1:])
	if argsiface != nil {
		*argsiface.(*[]string) = make([]string, len(flags.Args()))
		copy(*argsiface.(*[]string), flags.Args())
	}

	flags.Visit(noteVisited) // note all the visited flags, needed below

	// Set all pointer-type flags that actually had values set
	for kk := range pointers {
		if visited[kk] {
			pp, _ := sstype.FieldByName(kk)
			vv := ssvalue.FieldByName(kk)
			switch pp.Type.String() {
			case "*string":
				vv.Set(reflect.ValueOf(pointers[pp.Name]))
			case "*int":
				vv.Set(reflect.ValueOf(pointers[pp.Name]))
			case "*bool":
				vv.Set(reflect.ValueOf(pointers[pp.Name]))
			case "*int64":
				vv.Set(reflect.ValueOf(pointers[pp.Name]))
			case "*float64":
				vv.Set(reflect.ValueOf(pointers[pp.Name]))
			default:
				continue
			}
		}
	}
}
