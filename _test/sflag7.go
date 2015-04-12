package main

import (
	"fmt"
	"github.com/LDCS/sflag"
)

var opt = struct {
	SomeCommand string "! is command that might contain pipe char ! 'yes | head'"
}{}

// illustrates how to override the "Default value" delineator.  Provide your replacement as the first non-alphabetic of the tag
func main() {
	sflag.Parse(&opt)
	fmt.Println("SomeCommand=", opt.SomeCommand)
}
