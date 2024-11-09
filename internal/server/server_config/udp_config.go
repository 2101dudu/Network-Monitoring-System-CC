package server_config

import (
	"fmt"
	"net"
	p "nms/pkg/packet"
	u "nms/pkg/utils"
	"os"
)

var mapOfAgents map[byte]bool

func StartUDPServer(port string) {
	// Initialize the map
	mapOfAgents = make(map[byte]bool)

	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to resolve address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to initialize the server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[UDP] Server listening on port", port)

	for {
		handleUDPConnection(conn)
	}
}

func handleUDPConnection(conn *net.UDPConn) {
	fmt.Println("[UDP] Waiting for data from an agent")

	// Read registration request
	n, udpAddr, responseData := u.ReadUDP(conn, "[UDP] Registration request received", "[UDP] [ERROR] Unable to read registration request")

	// Check if there is data
	if n == 0 {
		fmt.Println("[UDP] [ERROR] No data received")
		return
	}

	// Check message type
	msgType := u.MessageType(responseData[0])
	switch msgType {
	case u.ACK:
		fmt.Println("[UDP] Acknowledgement received")
		return
	case u.METRICSGATHERING:
		fmt.Println("[UDP] Metrics received")
		return
	case u.REGISTRATION:
		// CHANGE TO THREAD
		fmt.Println("[UDP] Processing registration request...")

		// Decode registration request
		reg, err := p.DecodeRegistration(responseData[1:n])
		if err != nil {
			fmt.Println("[UDP] [ERROR] Unable to decode registration data:", err)
			// ****** SEND NOACK ******
			return
		}

		// Register agent
		mapOfAgents[reg.AgentID] = true

		// ****** SEND ACK ******
		p.SendAck(conn, udpAddr, reg.PacketID, reg.AgentID, true)
		return
	default:
		fmt.Println("[UDP] [ERROR] Unknown message type")
		return
	}
}
