package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	Mu              sync.Mutex
	clients         = make(map[string]client)
	NumberOfClients = 0
	messages        = make(chan message)
	OldMessages     = []string{}
)

type client struct {
	Name string
	Conn net.Conn
}

type message struct {
	Name    string
	Text    string
	Address string
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	content, err := os.OpenFile("linuxlogo.txt", os.O_CREATE, 0o755)
	if err != nil {
		log.Println(err)
	}
	logo, err := io.ReadAll(content)
	if err != nil {
		log.Println(err)
	}
	_, err = fmt.Fprintf(conn, "Welcome to TCP-Chat!\n%s", logo)
	if err != nil {
		deleteUserWithoutName(conn)
		return
	}
	_, err = fmt.Fprintf(conn, "\n[ENTER YOUR NAME]:")
	if err != nil {
		deleteUserWithoutName(conn)
		return
	}

	name := GetName(conn, 0)
	GiveHistory(conn)

	Mu.Lock()
	clients[conn.RemoteAddr().String()] = client{name, conn}
	Mu.Unlock()

	messages <- NewMessage(name+" has joined our server", conn)
	fmt.Fprintf(conn, CurrentTime()+"["+name+"]:")
	input := bufio.NewScanner(conn)
	for input.Scan() {
		_, err = fmt.Fprintf(conn, CurrentTime()+"["+name+"]:")
		if err != nil {
			deleteUser(name, conn)
			return
		}
		Mu.Lock()
		if IsValidString(input.Text()) && IsValidTxt(input.Text()) {
			messages <- NewMessage(CurrentTime()+"["+name+"]"+input.Text(), conn)
		}
		Mu.Unlock()
	}
	deleteUser(name, conn)
}

func deleteUserWithoutName(conn net.Conn) {
	Mu.Lock()
	NumberOfClients--
	delete(clients, conn.RemoteAddr().String())
	Mu.Unlock()
}

func deleteUser(name string, conn net.Conn) {
	messages <- NewMessage(name+" left our server", conn)
	Mu.Lock()
	NumberOfClients--
	delete(clients, conn.RemoteAddr().String())
	Mu.Unlock()
}

func IsValidTxt(text string) bool {
	var count int
	for _, i := range text {
		if i == '\n' || i == ' ' || i == '\t' {
			count++
		}
	}
	if count == len(text) {
		return false
	}
	return true
}

func NewMessage(text string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		Text:    text,
		Address: addr,
	}
}

func CurrentTime() string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "]"
}

func GetName(conn net.Conn, flag int) string {
	if flag == 1 {
		fmt.Fprintf(conn, "\n[ENTER YOUR NAME]:")
	}
	namereader := bufio.NewReader(conn)
	name, err := namereader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		deleteUserWithoutName(conn)
	}
	name = strings.TrimSpace(name)
	if !IsVacantName(name) {
		_, err := fmt.Fprintf(conn, "that name is taken, try again")
		if err != nil {
			deleteUserWithoutName(conn)
		}
		return GetName(conn, 1)
	} else if !IsValidString(name) {
		_, err := fmt.Fprintf(conn, "that name is invalid, try again")
		if err != nil {
			deleteUserWithoutName(conn)
		}
		return GetName(conn, 1)
	} else if len(name) > 10 {
		_, err := fmt.Fprintf(conn, "that name is too long, try again")
		if err != nil {
			deleteUserWithoutName(conn)
		}
		return GetName(conn, 1)
	}
	return name
}

func IsVacantName(name string) bool {
	Mu.Lock()
	for _, client := range clients {
		if strings.ToLower(name) == strings.ToLower(client.Name) {
			Mu.Unlock()
			return false
		}
	}
	Mu.Unlock()
	return true
}

func IsValidString(name string) bool {
	name = strings.TrimSuffix(name, "\n")
	rxmsg := regexp.MustCompile("^[\u0400-\u04FF\u0020-\u007F]+$")
	if !rxmsg.MatchString(name) {
		return false
	}

	return true
}
