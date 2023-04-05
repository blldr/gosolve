package main

import (
	"fmt"

)

func main() {
	roots, err := FindRoots("y + 6", -10, 0, 0.1, 0.01)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(roots)
}
