package main

import (
	"fmt"
	"net"
	"os"
)

const PORT = "8080"

func main() {
	fmt.Println("Starting Listener!")

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", PORT))
	if err != nil {
		fmt.Printf("Failed to bind to port %s\n", PORT)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Printf("Server is listening on port %s\n", PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 128)
	_, err := conn.Read(b)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Message received: ", b)
	_, err = conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println("Error writing connection: ", err.Error())
		os.Exit(1)
	}
}
