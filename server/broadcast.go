package server

import (
	"fmt"
	"net"
)

func GiveHistory(conn net.Conn) {
	Mu.Lock()
	for _, msg := range OldMessages {
		fmt.Fprintf(conn, msg+"\n")
	}
	Mu.Unlock()
}

func Broadcaster() {
	for {
		select {
		case msg := <-messages:
			OldMessages = append(OldMessages, msg.Text)
			Mu.Lock()
			for _, client := range clients {
				if msg.Address == client.Conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintf(client.Conn, "\n\033[1A\033[K"+msg.Text+"\n"+CurrentTime()+"["+client.Name+"]")

			}
			Mu.Unlock()
		}
	}
}
