package main

import (
	"fmt"
	"net"
	"os"
	"sync"
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

	var wg sync.WaitGroup
	for {
		conn, err := l.Accept()
		fmt.Printf("Received a connection, address:%s\n", conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		wg.Add(1)
		go handleClient(conn, &wg)
	}

	wg.Wait()
}

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer func() {
		fmt.Printf("Client Exit!\n")
		wg.Done()
		conn.Close()
	}()
	b := make([]byte, 1024)
	_, err := conn.Read(b)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Connection data, address:%s, message:%s\n", conn.RemoteAddr().String(), string(b))
	_, err = conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println("Error writing connection: ", err.Error())
		os.Exit(1)
	}
}
