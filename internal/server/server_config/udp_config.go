package server_config

import (
	"fmt"
	"net"
	m "nms/pkg/message"
	u "nms/pkg/utils"
	"os"
)

var agentCounter byte = 1

func StartUDPServer(port string) {
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
		agentCounter++
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
	case u.REGISTRATION:
		// CHANGE TO THREAD
		fmt.Println("[UDP] Processing registration request...")

		// Decode registration request
		reg, err := m.DecodeRegistration(responseData[1:n])
		if err != nil {
			fmt.Println("[UDP] [ERROR] Unable to decode registration data:", err)
			os.Exit(1)
		}

		// Validate registration request
		if reg.NewID != 0 || reg.SenderIsServer {
			fmt.Println("[UDP] [ERROR] Invalid registration request parameters")
			// ****** SEND NOACK ******
			return
		}

		// Create new registration request
		newReg := m.NewRegistrationBuilder().IsServer().SetNewID(agentCounter).Build()

		// Encode new registration request
		newRegData := m.EncodeRegistration(newReg)

		// Send new registration request
		u.WriteUDP(conn, udpAddr, newRegData, "[UDP] New registration request sent", "[UDP] [ERROR] Unable to send new registration request")

		// ****** SEND ACK ******

	case u.ACK:
		fmt.Println("[UDP] Acknowledgement received")

	case u.ERROR:
		fmt.Println("[UDP] Error message received")

	default:
		fmt.Println("[UDP] [ERROR] Unknown message type")

	}
}
