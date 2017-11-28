/*
There is a pattern to our pipeline functions:

1. stages close their outbound channels when all the send operations are done.
2. stages keep receiving values from inbound channels until those channels are closed
*/
package main

import (
	"fmt"
	"sync"
)

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

	//Distribute the sq work across two goroutines that both read from in.
	//This called fan-out
	var c1 = sq(c)
	var c2 = sq(c)

	//This called fan-in
	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}

//Send on closed channel panic, so it's important to ensure all send are done before calling close.
//The sync.WaitGroup provide the simple way to arrange this synchronization.
func merge(cs ...<-chan int) chan int {
	var wg sync.WaitGroup
	var out = make(chan int)
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
