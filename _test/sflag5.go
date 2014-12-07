package main

import (
	"github.com/LDCS/sflag"
	"fmt"
)

var opt = struct {
	Foo	*int	"foo"
}{}

func main() {
	foo	:= 1
	opt.Foo	= &foo	// will be ignored by sflag because it is a non-nil pointer
	sflag.Parse(&opt)
	fmt.Println("Foo", *opt.Foo)
}
