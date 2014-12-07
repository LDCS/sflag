package main

import (
	"github.com/LDCS/sflag"
	"fmt"
)

var opt = struct {
	Iq          int     "Upcased - winuser              | 42"	// Provide default to the right of Pipe, whitespace is trimmed
	Iq_         int     "Downcase - linuser             | 142"	// Underscore to the right of member will parse for --iq instead of --Iq
}{}

func main() {
	sflag.Parse(&opt)
	fmt.Println("Iq=", opt.Iq)
	fmt.Println("iq=", opt.Iq_)
}
