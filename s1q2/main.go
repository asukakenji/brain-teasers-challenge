package main

import (
	"./lib"
	"fmt"
)

func main() {
	digits := int(123456789)
	sum := int(100)
	fmt.Println("start", digits)
	channel := lib.Compute(sum, digits)
	for expression := range channel {
		fmt.Printf("%s = %d\n", expression, sum)
	}
	fmt.Println(digits)
}
