package agent_config

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func ConnectUDP(serverAddr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("[ERROR] Unable to resolve UDP address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("[ERROR] Unable to connect via UDP:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Write a UDP message: ")
		message, _ := reader.ReadString('\n')
		conn.Write([]byte(message))

		buffer := make([]byte, 1024)
		n, _, _ := conn.ReadFromUDP(buffer)
		fmt.Printf("UDP response from the server: %s\n", string(buffer[:n]))
	}
}
