package main

import (
	"github.com/LDCS/sflag"
	"fmt"
)

var opt = struct {
	Args	[]string
}{}

func main() {
	sflag.Parse(&opt)
	for ii, aa := range opt.Args {
		fmt.Println("arg num", ii, ":", aa)
	}
}
