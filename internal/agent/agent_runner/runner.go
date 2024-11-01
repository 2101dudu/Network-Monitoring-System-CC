package main

import (
	"fmt"
	ac "nms/internal/agent/agent_config"
	"os"
)

func main() {
	fmt.Print("Choose connection protocol (tcp/udp): ")
	var protocol string
	fmt.Scanln(&protocol)

	fmt.Print("Server address to connect to (e.g., localhost:8080): ")
	var serverAddr string
	fmt.Scanln(&serverAddr)

	if protocol == "tcp" {
		ac.ConnectTCP(serverAddr)
	} else if protocol == "udp" {
		ac.ConnectUDP(serverAddr)
	} else {
		fmt.Println("[URROR] Unknown procotol")
		os.Exit(1)
	}
}
