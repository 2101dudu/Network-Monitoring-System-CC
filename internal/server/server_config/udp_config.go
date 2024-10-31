package server_config

import (
	"fmt"
	"net"
	"strings"
)

func StartUDPServer(port string) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Println("[ERROR] Unable to resolve UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[ERROR] Unable to initialize UDP server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server listening on port", port)

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("[ERROR] Unable to read UDP data:", err)
			continue
		}
		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("UDP message received: %s\n", message)
		conn.WriteToUDP([]byte(message), clientAddr)
	}
}
