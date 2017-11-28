package main

import "fmt"

//The first stage gen is a function that converts a list of integers to a channel that emits the integers in the list.
func gen(nums ...int) <-chan int {
	var out = make(chan int)
	go func() {
		for _, num := range nums {
			out <- num
		}
		close(out)
	}()
	return out
}

//The second stage sq is a function that receives integers from a channel and returns a channel that emits the square of each received integer.
func sq(in <-chan int) <-chan int {
	var out = make(chan int)
	go func() {
		for num := range in {
			out <- num * num
		}
		close(out)
	}()

	return out
}

//The main function setup the pipeline and runs the final stage: it receives values from the second stage, and prints each one, until the channel is closed.
func main() {
	var c = gen(7, 9)
	var out = sq(c)

	for n := range out {
		fmt.Println(n)
	}
}
