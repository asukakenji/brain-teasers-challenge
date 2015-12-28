package main

import (
	"./lib"
	"bytes"
	"fmt"
)

func String(arr []int) string {
	if len(arr) == 0 {
		return "[]"
	}
	var buffer bytes.Buffer
	sep := "[ "
	for _, elem := range arr {
		elemString := fmt.Sprintf("%d", elem)
		buffer.WriteString(sep)
		buffer.WriteString(elemString)
		sep = ", "
	}
	buffer.WriteString(" ]")
	return buffer.String()
}

func main() {
	numbers := []int{7, 7, 4, 0, 9, 8, 2, 4, 1, 9}
	fmt.Println("start", String(numbers))
	lib.Partition(numbers)
	fmt.Println(String(numbers))
}
