package udp

import (
	"fmt"
	"net"
	"os"
)

func StartUDPServer(port string) {
	// Initialize the map
	agentsIPs = make(map[byte][4]byte)

	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Println("[SERVER] [ERROR 8] Unable to resolve address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[SERVER] [ERROR 9] Unable to initialize the server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[SERVER] Server listening on port", port)

	for {
		handleUDPConnection(conn)
	}
}
