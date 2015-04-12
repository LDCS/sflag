package main

import (
	"fmt"
	"github.com/LDCS/sflag"
)

var opt = struct {
	Age  int "age | 42"
	Args []string
}{
	Args: []string{"--Age", "10", "hello", "world"}, // One way to set what sflag will parse instead of os.Args[1:]
}

// shows to set the []string to be parsed instead of os.Args[1:]
func main() {
	sflag.Parse(&opt)
	fmt.Println("Age", opt.Age)
	for ii, aa := range opt.Args {
		fmt.Println("arg num", ii, ":", aa)
	}
}
