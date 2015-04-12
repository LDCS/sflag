package main

import (
	"fmt"
	"github.com/LDCS/sflag"
)

var opt = struct {
	Usage string "demonstrates upcase and lowcase flags"
	Iq    int    "Upcased - winuser              | 42"  // Provide default to the right of Pipe, whitespace is trimmed
	Iq_   int    "Downcase - linuser             | 142" // Underscore to the right of member will parse for --iq instead of --Iq
}{}

// illustrates how sflag sets the Usage member to the usage string composed from flag members.
func main() {
	fmt.Printf("Usage before sflag.Parse(%s)\n", opt.Usage)
	sflag.Parse(&opt)
	fmt.Printf("Usage after sflag.Parse(%s)\n", opt.Usage)
}
