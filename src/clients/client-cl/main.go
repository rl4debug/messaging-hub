package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
)

func main() {
	var red = color.New(color.FgRed)
	var green = color.New(color.FgGreen)

	fmt.Print("Enter a name: ")

	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')

	var serverPort = 9679
	var serverAddr = fmt.Sprintf("%s:%d", "localhost", serverPort)
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}

	var nameBytes = []byte(name)
	nameBytes = append([]byte{byte(1)}, nameBytes...)
	nameBytes = append([]byte{byte(len(nameBytes))}, nameBytes...)

	//Send name to server
	_, err = conn.Write(nameBytes)
	if err != nil {
		panic(err)
	}

	green.Println("Connected. Type something and Enter.")

	//Start send and receive message
	go func() {
		for {
			text, _ := reader.ReadString('\n')
			var data = []byte{byte(2)}
			data = append(data, []byte(text)...)
			data = append([]byte{byte(len(data))}, data...)
			conn.Write(data)
		}
	}()

	go func() {
		var mainData = make([]byte, 0)
		for {
			var data = make([]byte, 1)
			n, err := conn.Read(data)
			if err != nil {
				panic(err)
			}

			data = data[:n]
			mainData = append(mainData, data...)
			var l = int(mainData[0])
			if l < len(mainData) {
				var messageByte = mainData[1 : l+1]
				//ignore msgType
				messageByte = messageByte[1:]

				fromLen := int(messageByte[0])
				from := string(messageByte[1:fromLen])
				messageByte = messageByte[fromLen+1:]
				if messageByte[len(messageByte)-1] == 10 {
					messageByte = messageByte[:len(messageByte)-1]
				}
				text := string(messageByte)
				red.Print(from, ": ")
				fmt.Println(text)

				mainData = mainData[l+1:]
			}
		}
	}()

	var c = make(chan struct{})
	<-c //Wait
}
