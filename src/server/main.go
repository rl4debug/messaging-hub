package main

import (
	"fmt"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var hub = MakeHub()
	//Default port
	var port = 9679
	var err = hub.Listen(port)

	if err != nil {
		panic(err)
	}

	fmt.Println("Listening on port ", port)

	var c = make(chan struct{})
	<-c //Wait
}
