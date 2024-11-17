package udp

import (
	"fmt"
	"net"
	"os"
)

func getUDPConnection(serverAddr string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("[AGENT] [ERROR 1] Unable to resolve address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("[AGENT] [ERROR 2] Unable to connect:", err)
		os.Exit(1)
	}
	return conn
}
