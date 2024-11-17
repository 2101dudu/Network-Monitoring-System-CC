package udp

import (
	"fmt"
	"net"
	packet "nms/pkg/packet"
	utils "nms/pkg/utils"
	"os"
	"sync"
)

var (
	packetsWaitingAck = make(map[byte]bool)
	pMutex            sync.Mutex
)

func HandleUDPConnection(conn *net.UDPConn) {
	// generate Agent ID
	agentID, err := utils.GetAgentID()
	if err != nil {
		fmt.Println("[AGENT] [ERROR 3] Unable to get agent ID:", err)
		os.Exit(1)
	}
	// create registration request
	reg := packet.NewRegistrationBuilder().SetPacketID(1).SetAgentID(agentID).Build()
	// encode registration request
	regData := packet.EncodeRegistration(reg)

	// set the status of the packet to "not" waiting for ack
	packet.PacketIDIsWaiting(reg.PacketID, packetsWaitingAck, &pMutex, false)

	successMessage := "[AGENT] Registration request sent"
	errorMessage := "[AGENT] [ERROR 4] Unable to send registration request"
	go packet.SendPacketAndWaitForAck(reg.PacketID, packetsWaitingAck, &pMutex, conn, nil, regData, successMessage, errorMessage)

	for {
		fmt.Println("[AGENT] [MAIN READ THREAD] Waiting for response from server")

		// read message from server
		n, _, responseData := utils.ReadUDP(conn, "[AGENT] [MAIN READ THREAD] Response received", "[AGENT] [MAIN READ THREAD] [ERROR 5] Unable to read response")

		// Check if data was received
		if n == 0 {
			fmt.Println("[AGENT] [MAIN READ THREAD] [ERROR 6] No data received")
			return
		}

		// Check message type
		msgType := utils.MessageType(responseData[0])
		msgPayload := responseData[1:n]

		go func() {
			switch msgType {
			case utils.ACK:
				packet.HandleAck(msgPayload, packetsWaitingAck, &pMutex, agentID)
				return
			case utils.TASK:
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
