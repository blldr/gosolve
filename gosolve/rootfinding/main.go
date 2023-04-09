package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		os.Exit(1)
	}
	roots, err := FindRoots(args[1], -5, 5, 'y')
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("roots:",roots)
}
