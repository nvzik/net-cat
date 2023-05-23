package main

import (
	"netcat/server"
	"os"
)

func main() {
	server.RunServer(os.Args)
}
