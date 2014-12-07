package main

import (
	"github.com/LDCS/sflag"
)

var opt = struct {
	}{}

func main() {
	sflag.Parse(&opt)
}
