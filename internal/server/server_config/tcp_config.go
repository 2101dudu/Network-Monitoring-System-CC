package server_config

import (
	"bufio"
	"fmt"
	"net"
)

func StartTCPServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("[ERROR] Unable to start TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[ERROR] unable to accept TCP connection:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[ERROR] Unable to read TCP data:", err)
			break
		}
		fmt.Printf("TCP Message received: %s", message)
		conn.Write([]byte(message))
	}
}
