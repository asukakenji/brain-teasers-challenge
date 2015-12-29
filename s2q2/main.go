package main

import (
	"./lib"
)

func main() {
	queue := lib.NewQueue()
	queue.Add("Hey")
	queue.Add("there")
	queue.Add("world.")
	queue.Add("How")
	queue.Add("are")
	queue.Add("you?")
}
