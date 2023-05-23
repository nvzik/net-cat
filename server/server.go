package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func RunServer(args []string) {
	var port string
	if len(args) == 1 {
		port = "8989"
	}
	if len(args) == 2 {
		if portChecker(args[1]) {
			port = args[1]
		} else {
			log.Println("Invalid port.")
			return
		}
	}
	if len(args) > 2 {
		fmt.Println("[USAGE]: go run . $port")
		return
	}
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("The server is runnin' on localhost:" + port)
	defer listener.Close()
	go Broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		Mu.Lock()
		if NumberOfClients > 9 {
			fmt.Fprintf(conn, "server is full, try again later")
			conn.Close()
		} else {
			NumberOfClients++
			go HandleConn(conn)
		}
		Mu.Unlock()
	}
}

func portChecker(arg string) bool {
	port, err := strconv.Atoi(arg)
	if err != nil {
		return false
	}
	if port < 1024 || port > 65535 {
		return false
	}
	return true
}
