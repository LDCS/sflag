package main

import (
	"github.com/LDCS/sflag"
	"fmt"
)

var opt = struct {
	Usage		string	"demonstrates upcase and lowcase flags"
	Iq              int     "Upcased - winuser              | 42"	// Provide default to the right of Pipe, whitespace is trimmed
	Iq_             int     "Downcase - linuser             | 142"	// Underscore to the right of member will parse for --iq instead of --Iq
}{}

func main() {
	fmt.Printf("Usage before sflag.Parse(%s)\n", opt.Usage)
	sflag.Parse(&opt)
	fmt.Printf("Usage after sflag.Parse(%s)\n", opt.Usage)
}
