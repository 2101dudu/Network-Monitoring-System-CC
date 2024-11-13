package agent_config

import (
	"fmt"
	"net"
	p "nms/pkg/packet"
	u "nms/pkg/utils"
	"os"
	"sync"
)

func ConnectUDP(serverAddr string) {
	conn := getUDPConnection(serverAddr)

	defer conn.Close()

	handleUDPConnection(conn)
}

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

func handleUDPConnection(conn *net.UDPConn) {
	var (
		packetsWaitingAck = make(map[byte]bool)
		pMutex            sync.Mutex
	)

	// generate Agent ID
	agentID, err := u.GetAgentID()
	if err != nil {
		fmt.Println("[AGENT] [ERROR 3] Unable to get agent ID:", err)
		os.Exit(1)
	}
	// create registration request
	reg := p.NewRegistrationBuilder().SetPacketID(1).SetAgentID(agentID).Build()
	// encode registration request
	regData := p.EncodeRegistration(reg)

	pMutex.Lock()
	packetsWaitingAck[reg.PacketID] = false
	pMutex.Unlock()

	successMessage := "[AGENT] Registration request sent"
	errorMessage := "[AGENT] [ERROR 4] Unable to send registration request"
	go p.SendPacketAndWaitForAck(reg.PacketID, regData, packetsWaitingAck, &pMutex, conn, successMessage, errorMessage)

	for {
		fmt.Println("[AGENT] [MAIN READ THREAD] Waiting for response from server")

		// read message from server
		n, _, responseData := u.ReadUDP(conn, "[AGENT] [MAIN READ THREAD] Response received", "[AGENT] [MAIN READ THREAD] [ERROR 5] Unable to read response")

		// Check if data was received
		if n == 0 {
			fmt.Println("[AGENT] [MAIN READ THREAD] [ERROR 6] No data received")
			return
		}

		// Check message type
		msgType := u.MessageType(responseData[0])
		msgPayload := responseData[1:n]

		go func() {
			switch msgType {
			case u.ACK:
				p.HandleAck(msgPayload, packetsWaitingAck, &pMutex, agentID)
				return
			case u.TASK:
				fmt.Println("[AGENT] Metrics received from server")
				// HandleTask method - TO DO
				return
			default:
				fmt.Println("[AGENT] [ERROR 7] Unknown message type received from server")
				return
			}
		}()
	}
}
