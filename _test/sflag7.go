package main

import (
	"github.com/LDCS/sflag"
	"fmt"
)

var opt = struct {
		SomeCommand string  "! is command that might contain pipe char ! 'yes | head'"
	}{}

func main() {
	sflag.Parse(&opt)
	fmt.Println("SomeCommand=", opt.SomeCommand)
}
