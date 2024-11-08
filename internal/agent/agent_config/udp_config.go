package agent_config

import (
	"fmt"
	"net"
	p "nms/pkg/packet"
	u "nms/pkg/utils"
	"os"
)

func ConnectUDP(serverAddr string) {
	conn := getUDPConnection(serverAddr)

	defer conn.Close()

	handleUDPConnection(conn)
}

func getUDPConnection(serverAddr string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to resolve address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to connect:", err)
		os.Exit(1)
	}
	return conn
}

func handleUDPConnection(conn *net.UDPConn) {
	packetsWaitingAck := make(map[byte]bool)

	// generate Agent ID
	agentID, err := u.GetAgentID()
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to get agent ID:", err)
		os.Exit(1)
	}
	// create registration request
	reg := p.NewRegistrationBuilder().SetPacketID(1).SetAgentID(agentID).Build()
	// encode registration request
	regData := p.EncodeRegistration(reg)

	packetsWaitingAck[reg.PacketID] = false
	go func() {
		for {
			waiting, exists := packetsWaitingAck[reg.PacketID]

			if !exists { // registration packet has been removed from map
				break
			}
			if !waiting { // [CAUTION] TIMEOUT REQUIRED
				u.WriteUDP(conn, nil, regData, "[UDP] Registration request sent", "[UDP] [ERROR] Unable to send registration request")
				packetsWaitingAck[reg.PacketID] = true
			}
		}
	}()

	for {
		fmt.Println("[UDP] Waiting for response from server")

		// read message from server
		n, _, responseData := u.ReadUDP(conn, "[UDP] Response received", "[UDP] [ERROR] Unable to read response")

		// Check if data is received
		if n == 0 {
			fmt.Println("[UDP] [ERROR] No data received")
			return
		}

		// Check message type
		msgType := u.MessageType(responseData[0])
		msgPayload := responseData[1:n]

		go func() {
			switch msgType {
			case u.ACK:
				p.HandleAck(msgPayload, packetsWaitingAck, agentID)
				return
			case u.TASK:
				fmt.Println("[UDP] Metrics received from server")
				// HandleTask method - TO DO
				return
			default:
				fmt.Println("[UDP] [ERROR] Unknown message type received from server")
				return
			}
		}()
	}
}
