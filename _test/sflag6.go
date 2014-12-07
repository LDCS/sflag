package main

import (
	"github.com/LDCS/sflag"
	"fmt"
)

var opt = struct {
	Bar	*int	"bar"	// Bar will remain nil if flag was not not set.  Note default makes no sense here
}{}

func main() {
	sflag.Parse(&opt)
	if opt.Bar != nil {
		fmt.Println("Bar=", *opt.Bar)
	} else {
		fmt.Println("Bar was not set")
	}
}
