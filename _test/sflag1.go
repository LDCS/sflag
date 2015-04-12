package main

import (
	"github.com/LDCS/sflag"
)

var opt = struct {
}{}

// Minimal example
func main() {
	sflag.Parse(&opt)
}
