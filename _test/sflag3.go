package main

import (
	"fmt"
	"github.com/LDCS/sflag"
)

var opt = struct {
	Args []string
}{}

// shows how to retrieve unparsed flags out of os.Args[1:]
func main() {
	sflag.Parse(&opt)
	for ii, aa := range opt.Args {
		fmt.Println("arg num", ii, ":", aa)
	}
}
