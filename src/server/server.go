package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"

	"github.com/rl4debug/messaging-hub/src/message"
	"github.com/rl4debug/messaging-hub/src/simple-blockchain"
	"github.com/satori/go.uuid"
)

//Client hehe
type Client struct {
	Id          string
	Name        string
	Conn        net.Conn
	Hub         *Hub
	LastMessage message.Message
}

func (cli *Client) Start() {
	go func() {
		var processData = make([]byte, 0)
		reader := bufio.NewReader(cli.Conn)
		for {
			data := make([]byte, 1)
			n, err := reader.Read(data)
			if err != nil {
				//Connection closed, should remove this client
				if err == io.EOF {
					cli.Hub.Remove(cli)

					fmt.Println("remove client. Exit goroutine")
					//exit this goroutine
					return
				}
			}

			//Append to main buffer to process parse message
			processData = append(processData, data[:n]...)
			var l = int(processData[0])
			if l < len(processData) {
				var messageBytes = processData[1 : l+1]
				var msgType = int(messageBytes[0])
				switch msgType {
				case message.MSG_NAME:
					var msg = message.MessageRegister{
						Type: message.MSG_NAME,
						Name: string(messageBytes[1:]),
					}

					//Message is register name, so update name for client
					cli.Name = msg.Name
				case message.MSG_TEXT:
					msg := message.Message{
						Type: message.MSG_TEXT,
						Text: string(messageBytes[1:]),
					}

					//Talk to Hub to broadcast
					cli.LastMessage = msg

					//This maybe blocked
					cli.Hub.Notify <- cli
				default:
					fmt.Println("Invalid format")
					cli.Hub.Remove(cli)
					cli.Conn.Close()
					return
				}

				processData = processData[l+1:]
			}
		}
	}()
}

type Hub struct {
	Mu        sync.Mutex
	Clients   map[string]*Client
	Joins     chan net.Conn
	Notify    chan *Client
	ChainFile *os.File
	Chain     *blockchain.BlockChain
	Blocks    chan *blockchain.Block
}

func (hub *Hub) BroadcastExclude() {
	for client := range hub.Notify {
		//Update blockchain
		var block = hub.Chain.CreateNewBlock(client.LastMessage.Seriallize())
		hub.Blocks <- block

		for cli := range hub.LoopClients() {
			if cli.Id != client.Id {
				var msg = message.MessageBroadcast{
					Type: client.LastMessage.Type,
					From: client.Name,
					Text: client.LastMessage.Text,
				}
				messageBytes := msg.Seriallize()
				messageBytes = append([]byte{byte(len(messageBytes))}, messageBytes...)
				cli.Conn.Write(messageBytes)
			}
		}
	}
}

func (hub *Hub) LoopClients() <-chan *Client {
	var c = make(chan *Client)
	f := func() {
		hub.Mu.Lock()
		defer hub.Mu.Unlock()

		for _, client := range hub.Clients {
			c <- client
		}
		close(c)
	}

	go f()

	return c
}

func (hub *Hub) AddClient(client *Client) {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	hub.Clients[client.Id] = client
}

func (hub *Hub) Remove(client *Client) {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	delete(hub.Clients, client.Id)
}

func Listen(port int) (net.Listener, error) {
	var lstn net.Listener
	var err error
	var addr = fmt.Sprintf("%s:%d", "localhost", port)
	lstn, err = net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	return lstn, nil
}

func (hub *Hub) Listen(port int) error {
	var lstn, err = Listen(port)

	if err != nil {
		return err
	}

	//process Conn to be joined
	go hub.processJoins()

	go hub.BroadcastExclude()

	go hub.HandleNewBlock()

	//Accepts connections
	go func() {
		for {
			conn, err := lstn.Accept()
			if err != nil {
				fmt.Println("Accept connection error.")
				continue
			}

			//send connection to Joins queue
			hub.Joins <- conn
		}
	}()

	return nil
}

func (hub *Hub) HandleNewBlock() {
	for {
		select {
		case block := <-hub.Blocks:

			//Save block to file
			var data = block.Seriallize()
			var len = uint32(len(data))
			var lenBytes = make([]byte, 4)
			binary.LittleEndian.PutUint32(lenBytes, len)
			data = append(lenBytes, data...)
			hub.ChainFile.Write(data)
		}
	}
}

func (hub *Hub) processJoins() {
	for {
		select {
		case conn := <-hub.Joins:
			//TODO make client
			var client = &Client{
				Conn: conn,
				Id:   uuid.NewV1().String(),
				Hub:  hub,
			}
			hub.AddClient(client)
			client.Start()
		}
	}
}

func MakeHub() *Hub {
	ioutil.WriteFile("blockchaindata", nil, 0600)
	f, err := os.OpenFile("blockchaindata", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	var hub = &Hub{
		Joins:   make(chan net.Conn),
		Clients: make(map[string]*Client),
		Notify:  make(chan *Client),

		ChainFile: f,
		Chain:     &blockchain.BlockChain{},
		Blocks:    make(chan *blockchain.Block, 10),
	}

	//First block
	var firstBlock = hub.Chain.GenerateFirstBlock()
	hub.Chain.Blocks = append(hub.Chain.Blocks, firstBlock)
	hub.Blocks <- firstBlock

	return hub
}
