package server

import (
	"fmt"
	"net"
	packet "nms/pkg/packet"
	utils "nms/pkg/utils"
	"os"
)

var mapOfAgents map[byte]bool

func StartUDPServer(port string) {
	// Initialize the map
	mapOfAgents = make(map[byte]bool)

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

func handleUDPConnection(conn *net.UDPConn) {
	fmt.Println("[SERVER] [MAIN READ THREAD] Waiting for data from an agent")

	// Read registration request
	n, udpAddr, responseData := utils.ReadUDP(conn, "[SERVER] [MAIN READ THREAD] Registration request received", "[SERVER] [ERROR 10] Unable to read registration request")

	// Check if there is data
	if n == 0 {
		fmt.Println("[SERVER] [MAIN READ THREAD] [ERROR 11] No data received")
		return
	}

	// Check message type
	msgType := utils.MessageType(responseData[0])
	switch msgType {
	case utils.ACK:
		fmt.Println("[SERVER] Acknowledgement received")
		return
	case utils.METRICSGATHERING:
		fmt.Println("[SERVER] Metrics received")
		return
	case utils.REGISTRATION:
		// CHANGE TO THREAD
		fmt.Println("[SERVER] Processing registration request...")

		// Decode registration request
		reg, err := packet.DecodeRegistration(responseData[1:n])
		if err != nil {
			fmt.Println("[SERVER] [ERROR 12] Unable to decode registration data:", err)
			// ****** SEND NOACK ******
			return
		}

		// Register agent
		mapOfAgents[reg.AgentID] = true

		// ****** SEND ACK ******
		ack := packet.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).HasAcknowledged().Build()
		packet.EncodeAndSendAck(conn, udpAddr, ack)
		return
	default:
		fmt.Println("[SERVER] [ERROR 13] Unknown message type")
		return
	}
}
