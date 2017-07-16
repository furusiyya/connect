package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go Server(&wg)
	go Client(&wg)
	wg.Wait()

}

func Server(wg *sync.WaitGroup) {

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":8081")
	CheckErrors(err, "Server Listening!")

	// accept connection on port
	conn, err := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// will listen for messages from client
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Print("Client Disconnected!\n")
			break
		}

		if message == "exit\n" {
			conn.Close()
			break
		}

		// process the recieved string from client
		newmessage := strings.ToUpper(message)
		// send new string back to client
		_, err = conn.Write([]byte(newmessage + "\n"))

		if err != nil {
			fmt.Println("sendgin error!")
			panic(err)
		}

	}
	wg.Done()

}

func CheckErrors(err error, message string) {
	if err != nil {

		panic(err)
	}
	fmt.Println(message)

}

func Client(wg *sync.WaitGroup) {

	// connect to this socket
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	CheckErrors(err, "Client Connected!")
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Client: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text+"\n")
		// listen for reply
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Print("Disconnected\n")
			conn.Close()
			break
		}
		fmt.Print("Message from Server: " + message)
	}
	wg.Done()

}
