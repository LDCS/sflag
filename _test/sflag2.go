package main

import (
	"fmt"
	"github.com/LDCS/sflag"
)

var opt = struct {
	Iq  int "Upcased - winuser              | 42"  // Provide default to the right of Pipe, whitespace is trimmed
	Iq_ int "Downcase - linuser             | 142" // Underscore to the right of member will parse for --iq instead of --Iq
}{}

// shows how to parse for a simple int with upcased or lowcase first char of the flag
func main() {
	sflag.Parse(&opt)
	fmt.Println("Iq=", opt.Iq)
	fmt.Println("iq=", opt.Iq_)
}
