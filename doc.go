/*
Package sflag is a commandline arguments parser using struct syntax

Usage

	var opt = struct {
	    Usage       string  "sflags demonstrator"
	    SomeFile    string  "contains the something      | /dev/null"
	    IQ          int     "do not inflate              | 42"
	    GDP         float64 "in Vietnamese Dong          | 42000000000000000000000000.0"
	    Age         int64   "in milliseconds since epoch | 42000000000000"
	    SomeCommand string  "! is command that might contain pipe char ! 'yes | head'"
	    Verbose     bool    "Bool flags require use of an equals sign syntax (i.e. \"var=value\") to be unambiguous		| false"
	    OutData     string  " must be writable | /an/output/file"
	    Args	[]string
	    Foo	        *int			// sflag will ignore both member and tag since it is initialized to non-nil pointer
	    Bar	        *int    "bar"	        // sflag will not look for default value in this tag, since it will be a nil pointer
	    Baz_        int     "Set by --baz, not --Baz             | 42"
	    ignoreMe    string			// sflag will ignore low-case members
	}{}

	func Example() {
	     foo	:= 1;
	     opt.Foo	= &foo		// sflag will ignore this variable since it is a non-nil pointer

	     sflag.Parse(&opt)
	     fmt.Println("SomeFile=", opt.SomeFile)
	     fmt.Println("Age=", opt.Age)
	     fmt.Println("IQ=", opt.IQ)
	     fmt.Println("GDP=", opt.GDP)
	     fmt.Println("SomeCommand=", opt.SomeCommand)
	     fmt.Println("Verbose=", opt.Verbose)
	     if opt.Bar != nil {
	     	fmt.Println("Bar=", *opt.Bar)
	     } else {
		fmt.Println("Bar was not set")
	     }
	     fmt.Println("baz=", opt.Baz_)
	     for ii, aa := range opt.Args {
	     	fmt.Println("arg num", ii, ":", aa)
	     }
	}


*/
package sflag
