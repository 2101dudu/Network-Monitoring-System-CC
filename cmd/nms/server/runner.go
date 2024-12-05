package main

import (
	servertcp "nms/internal/server/alertflow"
	serverudp "nms/internal/server/nettask"
)

func main() {
	// Start tcp server
	go servertcp.StartTCPServer("8080")
	// start udp server
	serverudp.StartUDPServer("8081")
}
