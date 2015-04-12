package main

import (
	"fmt"
	"github.com/LDCS/sflag"
)

var opt = struct {
	Foo *int "foo"
}{}

// shows pointer members are ignored if non-nil
func main() {
	foo := 1
	opt.Foo = &foo // will be ignored by sflag because it is a non-nil pointer
	sflag.Parse(&opt)
	fmt.Println("Foo", *opt.Foo)
}
