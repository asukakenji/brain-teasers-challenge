package main

import (
	"./lib"
	"fmt"
)

func main() {
	digits := 123456789
	sum := 100
	fmt.Println("start", digits)
	channel := lib.Compute(digits)
	for expression := range channel {
		if expression.Value == sum {
			fmt.Printf("%s = %d\n", expression.Text, sum)
		}
	}
	fmt.Println(digits)
}
