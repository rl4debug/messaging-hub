# Simple Room chat include Server and Clients

Project structure (main items)
*   src/
    * server/ <- This is server
    * clients/
        * client-cl/ <- This is command-line client

---

### Note: (must install golang compiler first [Here](https://golang.org/dl/))

`1. Get source`

```sh
go get -u -v github.com/rl4debug/messaging-hub
#Command above will download source to `$GOPATH/src/github.com/rl4debug/message-hub`

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