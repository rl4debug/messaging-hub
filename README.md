# Simple Room chat include Server and Clients

Project structure (main items)
*   src/
    * server/ <- This is server
    * clients/
        * client-cl/ <- This is command-line client

---

### Note: (must install golang compiler [Here](https://golang.org/dl/), and GIT command-line)

`1. Get source`

```sh
#1. Create folder structure `$GOPATH/src/github.com/rl4debug/
#2. Navigate to `$GOPATH/src/github.com/rl4debug/ then run command:
git clone https://github.com/rl4debug/messaging-hub

#Download dependency packages
go get -u -v github.com/satori/go.uuid
go get -u -v github.com/fatih/color

```

`2. Build and Run`

```sh
#For build server
#navigate to ....github.com/rl4debug/message-hub/src/server
#Then run command:
go build main.go server.go

#Run server
./main


#Build client
#navigate to ....github.com/rl4debug/message-hub/src/clients/client-cl
#Then run command:
go build main.go
#Run client
./main

#Should run multiple clients to join the Hub
```