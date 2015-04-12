// Package sflag is a flag package variant that is 100% DRY, free of fugly pointer syntax and uses clean struct syntax.
//
// Implementation makes use of reflection and struct tags, in manner similar to previously published flag variants.
//
package sflag

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	visited map[string]bool
)

func noteVisited(_flag *flag.Flag) {
	visited[_flag.Name] = true
}

// Parse iterates through the members of the struct.  Notes:
//
//     Members are set up for std flag package to do the actual parsing, using type obtained via reflection and info from struct tag for usage and default setting.
//     Normally, the rightmost pipe char in the tag is used to delineate between Description (on left) and Default value (on right).
//     (You can override delineator to the first char of the tag (after eliminating leading whitespace) if such char is not alphabetic).
//     Fields with no tag or whitespace-only tags are ignored.
//     Non-nil pointer fields are ignored.
//     Nil pointer fields will be left nil if that flag is not set on commandline (and the tag is not parsed for a default value).
//     Flags starting with lowercase letter require that the coresponding member ends in single underscore.
//     Provide string member Usage initialized to brief program description.  Parse will append member descriptions to that string.
//     Provide []string member Args if you want to want to retrieve unconsumed flags.
//     Initialize []string member Args to the string array you want to parse instead of os.Args[1:].
func Parse(ss interface{}) { parseInternal(ss, true) }

// Parse2 is identical to Parse, except panics if there is both (1) a boolean flag and (2) a standalone true/false argument.
// It reminds you to use "--Foo=true" syntax (instead of "--Foo true" which would terminate the stdlib's flag processing for bool flag Foo, which is considered set by its presence alone).
// The downside of using this func is that unrelated presence of true/false results in progam panic.
func Parse2(ss interface{}) { parseInternal(ss, false) }

func parseInternal(ss interface{}, _permitStandaloneBool bool) {
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

	var argsiface interface{}
	args := make([]string, len(os.Args)-1)
	copy(args, os.Args[1:])

	progname := os.Args[0]
	if pp, ok := sstype.FieldByName("Args"); ok {
		if pp.Type.String() == "[]string" { // caller wanted to override os.Args and/or retrieve unconsumed flags
			vv := ssvalue.FieldByName("Args")
			if len(*vv.Addr().Interface().(*[]string)) == 0 {
			} else { // caller wanted to override os.Args
				args = make([]string, len(*vv.Addr().Interface().(*[]string)))
				copy(args, *vv.Addr().Interface().(*[]string))
			}
			argsiface = vv.Addr().Interface()
		}
	}

	moreusage := ""
	hasBoolArg := false
	flags := *flag.NewFlagSet(progname, flag.PanicOnError)

	for ii := 0; ii < sstype.NumField(); ii++ {
		pp := sstype.Field(ii)
		vv := ssvalue.Field(ii)
		switch {
		case pp.Anonymous:
			continue // Skip embedded fields
		case pp.Name == "Usage":
			continue // Not a flag
		case pp.Type.String() == "[]string":
			continue // Already handled Args, and not interested in other such members
		case (pp.Type.Kind() == reflect.Ptr) && (vv.Elem().Kind() != reflect.Invalid):
			continue // Ignore non-nil pointer members
		}

		tag := strings.TrimSpace((string)(pp.Tag))
		if tag == "" {
			continue
		}

		flagname := pp.Name
		if nn := len(pp.Name) - 1; flagname[nn] == '_' { // User wants to look for --f* instead of --F*
			flagname = strings.ToLower(pp.Name[:1]) + pp.Name[1:nn]
		}

		_, nn := utf8.DecodeRuneInString(tag)
		splitChar := tag[0:nn]
		if strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", splitChar) {
			splitChar = "|"
		} else {
			tag = tag[len(splitChar):]
		}

		lastSplit := strings.LastIndex(tag, splitChar)
		part0, part1 := "", ""
		switch lastSplit > -1 {
		case false:
			part1 = strings.TrimSpace(tag)
		case true:
			part0, part1 = strings.TrimSpace(tag[:lastSplit]), strings.TrimSpace(tag[(lastSplit+1):])
		}

		if pp.Type.Kind() == reflect.Ptr {
			switch pp.Type.String() {
			case "*string":
				tempstr := ""
				pointers[flagname] = &tempstr
				flags.StringVar(&tempstr, flagname, tempstr, "")
			case "*int":
				tempint := 0
				pointers[flagname] = &tempint
				flags.IntVar(&tempint, flagname, tempint, "")
			case "*bool":
				tempbool := false
				pointers[flagname] = &tempbool
				flags.BoolVar(&tempbool, flagname, tempbool, "")
			case "*int64":
				tempint64 := int64(0)
				pointers[flagname] = &tempint64
				flags.Int64Var(&tempint64, flagname, tempint64, "")
			case "*float64":
				tempfloat64 := 0.0
				pointers[flagname] = &tempfloat64
				flags.Float64Var(&tempfloat64, flagname, tempfloat64, "")
			default:
				continue
			}
		}

		if lastSplit < 0 {
			switch pp.Type.Kind() {
			case reflect.String:
				flags.StringVar(vv.Addr().Interface().(*string), flagname, vv.String(), " <--default, string # "+part0)
			case reflect.Int:
				flags.IntVar(vv.Addr().Interface().(*int), flagname, int(vv.Int()), " <--default, int # "+part0)
			case reflect.Bool:
				flags.BoolVar(vv.Addr().Interface().(*bool), flagname, bool(vv.Bool()), " <--default, bool # "+part0)
				hasBoolArg = true
			case reflect.Int64:
				flags.Int64Var(vv.Addr().Interface().(*int64), flagname, vv.Int(), " <--default, int64 # "+part0)
			case reflect.Float64:
				flags.Float64Var(vv.Addr().Interface().(*float64), flagname, vv.Float(), " <--default, float64 # "+part0)
			default:
				continue
			}
		}

		if lastSplit >= 0 {
			switch pp.Type.Kind() {
			case reflect.String:
				vv.SetString(part1)
				flags.StringVar(vv.Addr().Interface().(*string), flagname, part1, " <--default, string # "+part0)
			case reflect.Int:
				inum, _ := strconv.ParseInt(part1, 10, 64)
				vv.SetInt(inum)
				flags.IntVar(vv.Addr().Interface().(*int), flagname, int(inum), " <--default, int # "+part0)
			case reflect.Bool:
				bnum, _ := strconv.ParseBool(part1)
				vv.SetBool(bnum)
				flags.BoolVar(vv.Addr().Interface().(*bool), flagname, bool(bnum), " <--default, bool # "+part0)
				hasBoolArg = true
			case reflect.Int64:
				jnum, _ := strconv.ParseInt(part1, 10, 64)
				vv.SetInt(jnum)
				flags.Int64Var(vv.Addr().Interface().(*int64), flagname, jnum, " <--default, int64 # "+part0)
			case reflect.Float64:
				fnum, _ := strconv.ParseFloat(part1, 64)
				vv.SetFloat(fnum)
				flags.Float64Var(vv.Addr().Interface().(*float64), flagname, fnum, " <--default, float64 # "+part0)
			default:
				continue
			}
		}

		if lastSplit >= 0 {
			moreusage += "\n\t--" + flagname + ": " + part1 + " <-- Default, " + pp.Type.String() + " # " + part0
		}
	}

	if pp, ok := sstype.FieldByName("Usage"); ok {
		vv := ssvalue.FieldByName("Usage")
		vv.SetString("\n Usage of " + progname + " # " + (string)(pp.Tag) + "\n ARGS:" + moreusage)
	}

	if hasBoolArg && !_permitStandaloneBool {
		for _, arg := range args {
			switch strings.ToLower(arg) {
			case "true", "false":
				panic("Golang flag package requires \"--Foo=bar\" instead of \"--Foo bar\" syntax for bool args")
			}
		}
	}

	flags.Parse(args)
	if argsiface != nil {
		*argsiface.(*[]string) = make([]string, len(flags.Args()))
		copy(*argsiface.(*[]string), flags.Args())
	}

	flags.Visit(noteVisited) // note all the visited flags, needed below

	// Set all pointer-type flags that actually had values set
	for flagname := range pointers {
		if visited[flagname] {
			fieldname := flagname
			if flagname[:1] != strings.ToUpper(flagname[:1]) {
				fieldname = strings.ToUpper(flagname[:1]) + flagname[1:] + "_"
			}
			pp, _ := sstype.FieldByName(fieldname)
			vv := ssvalue.FieldByName(fieldname)
			switch pp.Type.String() {
			case "*string":
				vv.Set(reflect.ValueOf(pointers[flagname]))
			case "*int":
				vv.Set(reflect.ValueOf(pointers[flagname]))
			case "*bool":
				vv.Set(reflect.ValueOf(pointers[flagname]))
			case "*int64":
				vv.Set(reflect.ValueOf(pointers[flagname]))
			case "*float64":
				vv.Set(reflect.ValueOf(pointers[flagname]))
			default:
				continue
			}
		}
	}
}
