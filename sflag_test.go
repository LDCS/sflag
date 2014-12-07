package sflag

import (
	"fmt"
	"testing"
)

// TestSflag_1 is minimal
func TestSflag_1(*testing.T) {
	var opt = struct {
	}{}
	Parse(&opt)
}

// TestSflag_2 shows how to parse for a simple int with upcased or lowcase first char of the flag
func TestSflag_2(t *testing.T) {
	var opt = struct {
		Iq          int     "Upcased - winuser              | 42"	// Provide default to the right of Pipe, whitespace is trimmed
		Iq_         int     "Downcase - linuser             | 142"	// Underscore to the right of member will parse for --iq instead of --Iq
	}{}
	Parse(&opt)
	fmt.Println("Iq =", opt.Iq)
	fmt.Println("iq =", opt.Iq_)
	if opt.Iq != 42 || opt.Iq_ != 142 {
		t.Fail()
	}
	
}

// TestSflag_3 shows how to retrieve unparsed flags out of os.Args[1:]
func TestSflag_3(t *testing.T) {
	var opt = struct {
		Args	[]string
	}{}
	Parse(&opt)
	for ii, aa := range opt.Args {
		fmt.Println("arg num", ii, ":", aa)
	}
	if len(opt.Args) != 0 {
		t.Fail()
	}
}

// TestSflag_4 shows how to set the []string to be parsed instead of os.Args[1:]
func TestSflag_4(t *testing.T) {
	var opt = struct {
		Age     int      "Age | 20"
		Bar     int      "Bar | 150"
		GDP     int      "GDP | 10"
		Args	[]string
	}{
		Args : []string{"--Age", "10", "--Bar", "7", "--GDP", "2", "hello", "world"},	// One way to set what sflag will parse instead of os.Args[1:]
	}
	Parse(&opt)
	
	fmt.Println("Age =", opt.Age)
	fmt.Println("Bar =", opt.Bar)
	fmt.Println("GDP =", opt.GDP)
	
	for ii, aa := range opt.Args {
		fmt.Println("arg num", ii, ":", aa)
	}

	if opt.Age != 10 || opt.Bar != 7 || opt.GDP != 2 || len(opt.Args) != 2 || opt.Args[0] != "hello" || opt.Args[1] != "world" {
		t.Fail()
	}
}

// TestSflag_5 shows pointer members are ignored if non-nil
func TestSflag_5(t *testing.T) {
	var opt = struct {
		Foo	*int	"foo | 200"
	}{}
	foo	:= 1
	opt.Foo	= &foo
	Parse(&opt)
	fmt.Println("Foo", *opt.Foo)
	if foo != 1 {
		t.Fail()
	}
}

// TestSflag_6 shows nil pointer members remain nil if the corresponding flag was not set, else will point to the set value
func TestSflag_6(t *testing.T) {
	var opt1 = struct {
		Bar	*int	"bar"	// Bar will remain nil if flag was not not set.  Note default makes no sense here
	}{}
	Parse(&opt1)
	if opt1.Bar != nil {
		fmt.Println("Bar=", *opt1.Bar)
	} else {
		fmt.Println("Bar was not set")
	}

	var opt2 = struct {
		Bar	*int	"bar"	// Bar will remain nil if flag was not not set.  Note default makes no sense here
		Args    []string
	}{ Args : []string{"--Bar", "200"}}
	Parse(&opt2)
	if opt2.Bar != nil {
		fmt.Println("Bar=", *opt2.Bar)
	} else {
		fmt.Println("Bar was not set")
	}

	if opt1.Bar != nil || opt2.Bar == nil {
		t.Fail()
	}

	if *(opt2.Bar) != 200 {
		t.Fail()
	}

}

// TestSflag_7 illustrates how to override the "Default value" delineator.  Provide your replacement as the first non-alphabetic of the tag
func TestSflag_7(t *testing.T) {
	var opt = struct {
		SomeCommand string  "! is command that might contain pipe char ! 'yes | head'"
	}{}
	Parse(&opt)
	fmt.Println("SomeCommand=", opt.SomeCommand)

	if opt.SomeCommand != "'yes | head'" {
		t.Fail()
	}
}

// TestSflag_8 checks whether --Usage shows usage of the program
func TestSflag_8(t *testing.T) {
	var opt = struct {
		Usage	    string  "demonstrates upcase and lowcase flags"
		Iq          int     "Upcased - winuser              | 42"	// Provide default to the right of Pipe, whitespace is trimmed
		Iq_         int     "Downcase - linuser             | 142"	// Underscore to the right of member will parse for --iq instead of --Iq
	}{}
	var oldUsage = opt.Usage
	fmt.Printf("Usage before sflag.Parse(%s)\n", opt.Usage)
	Parse(&opt)
	fmt.Printf("Usage after sflag.Parse(%s)\n", opt.Usage)
	if oldUsage == opt.Usage {
		t.Fail()
	}
}

// TestSflag_9 is the kitchen sink demo
func TestSflag_9(t *testing.T) {
	var opt = struct {
		Usage       string  "sflags demonstrator"
		SomeFile    string  "contains the something      | /dev/null"
		IQ          int     "do not inflate              | 42"
		GDP         float64 "in Vietnamese Dong          | 42000000000000000000000000.0"
		Age         int64   "in milliseconds since epoch | 42000000000000"
		SomeCommand string  "! is command that might contain pipe char ! 'yes | head'"
		Verbose     bool    "Bool flags require use of an equals sign syntax (i.e. \"var=value\") to be unambiguous		| false"
		OutData     string  " must be writable | /an/output/file"
		Foo	    *int			// sflag will ignore both member and tag since it is initialized to non-nil pointer
		Bar	    *int    "bar"	        // sflag will not look for default value in this tag, since it will be a nil pointer
		Args	    []string
		ignoreMe    string			// sflag will ignore low-case members
	}{	Args : []string{"--Age", "10", "--Bar", "7", "--GDP", "2", "hello", "world"},	// One way to set what sflag will parse instead of os.Args[1:]
	}
	foo	:= 1;
	opt.Foo	= &foo		// sflag will ignore this variable since it is a non-nil pointer

	Parse(&opt)
	fmt.Println("SomeFile=", opt.SomeFile)
	fmt.Println("Age=", opt.Age)
	fmt.Println("IQ=", opt.IQ)
	fmt.Println("GDP=", opt.GDP)
	fmt.Println("SomeCommand=", opt.SomeCommand)
	fmt.Println("Verbose=", opt.Verbose)
	fmt.Println("OutData=", opt.OutData)
	if opt.Bar != nil {
		fmt.Println("Bar=", *opt.Bar)
	} else {
		fmt.Println("Bar was not set")
	}
	for ii, aa := range opt.Args {
		fmt.Println("arg num", ii, ":", aa)
	}

	if opt.SomeFile != "/dev/null" ||
		opt.Age != int64(10) ||
		opt.IQ != 42 ||
		opt.GDP != 2.0 ||
		opt.SomeCommand != "'yes | head'" ||
		opt.Verbose != false ||
		opt.OutData != "/an/output/file" ||
		opt.Bar == nil ||
		len(opt.Args) != 2 ||
		foo != 1 {
		t.Fail()
	}

	if opt.Args[0] != "hello" || opt.Args[1] != "world" || *(opt.Bar) != 7 {
		t.Fail()
	}
}

