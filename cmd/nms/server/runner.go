package main

import (
	servertcp "nms/internal/server/tcp"
	serverudp "nms/internal/server/udp"
)

func main() {
	// Start tcp server
	go servertcp.StartTCPServer("8080")
	// start udp server
	serverudp.StartUDPServer("8081")
}
