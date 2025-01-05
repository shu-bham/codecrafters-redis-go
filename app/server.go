package main

import (
	"log"
	"net"
	"os"
)

const PORT = "6379"

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:"+PORT)
	if err != nil {
		log.Printf("Failed to bind to port %s: %v\n", PORT, err)
		os.Exit(1)
	}
	defer l.Close()
	log.Printf("Server is listening on port %s\n", PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		log.Printf("Received a connection, address: %s\n", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("Error reading connection: %v\n", err)
			}
			return
		}

		if n == 0 {
			continue
		}

		log.Printf("Connection data, address: %s, message: %s\n",
			conn.RemoteAddr().String(), string(buf[:n]))

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			log.Printf("Error writing connection: %v\n", err)
			return
		}
	}
}
