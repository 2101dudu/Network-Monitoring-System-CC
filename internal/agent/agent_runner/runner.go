package main

import (
	"fmt"
	ac "nms/internal/agent/agent_config"
)

func main() {
	fmt.Print("Choose the protocol (tcp/udp): ")
	var protocol string
	fmt.Scanln(&protocol)

	fmt.Print("Write the server's address (e.g., localhost:8080): ")
	var serverAddr string
	fmt.Scanln(&serverAddr)

	if protocol == "tcp" {
		ac.ConnectTCP(serverAddr)
	} else if protocol == "udp" {
		ac.ConnectUDP(serverAddr)
	} else {
		fmt.Println("[ERROR] Unknown protocol")
	}
}
