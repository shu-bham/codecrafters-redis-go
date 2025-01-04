package main

import (
	"fmt"
	"net"
	"os"
)

const PORT = "6379"

func main() {
	fmt.Println("Starting Listener!")

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", PORT))
	if err != nil {
		fmt.Printf("Failed to bind to port %s\n", PORT)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Printf("Server is listening on port %s\n", PORT)

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	conn.Write([]byte("+PONG\r\n"))
}
