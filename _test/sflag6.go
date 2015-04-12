package main

import (
	"fmt"
	"github.com/LDCS/sflag"
)

var opt = struct {
	Bar *int "bar" // Bar will remain nil if flag was not not set.  Note default makes no sense here
}{}

// shows nil pointer members remain nil if the corresponding flag was not set, else will point to the set value
func main() {
	sflag.Parse(&opt)
	if opt.Bar != nil {
		fmt.Println("Bar=", *opt.Bar)
	} else {
		fmt.Println("Bar was not set")
	}
}
