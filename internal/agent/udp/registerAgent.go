package udp

import (
	"fmt"
	"net"
	ack "nms/internal/packet/ack"
	registration "nms/internal/packet/registration"
	utils "nms/internal/utils"
	"sync"
)

var (
	packetsWaitingAck = make(map[byte]bool)
	pMutex            sync.Mutex
)

var agentID byte

func registerAgent(conn *net.UDPConn, agentIP string) {
	var firstPacketID byte = 1
	var registrationData []byte
	agentID, registrationData = registration.CreateRegistrationPacket(firstPacketID, agentIP)

	// set the status of the packet to "not" waiting for ack, because it is yet to be sent
	utils.PacketIsWaiting(firstPacketID, packetsWaitingAck, &pMutex, false)

	successMessage := "[AGENT] Registration request sent"
	errorMessage := "[AGENT] [ERROR 4] Unable to send registration request"
	go ack.SendPacketAndWaitForAck(firstPacketID, packetsWaitingAck, &pMutex, conn, nil, registrationData, successMessage, errorMessage)

	ackWasSent := false
	for !ackWasSent {
		fmt.Println("[AGENT] [MAIN READ THREAD] Waiting for response from server")

		// read packet from server
		n, _, data := utils.ReadUDP(conn, "[AGENT] [MAIN READ THREAD] Response received", "[AGENT] [MAIN READ THREAD] [ERROR 5] Unable to read response")

		// Check if data was received
		if n == 0 {
			fmt.Println("[AGENT] [MAIN READ THREAD] [ERROR 6] No data received")
			return
		}

		// get ACK contents
		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		if packetType != utils.ACK {
			fmt.Println("[AGENT] [ERROR 17] Unexpected packet type received from server")
			return
		}
		ackWasSent = ack.HandleAck(packetPayload, packetsWaitingAck, &pMutex, agentID, conn)
	}
	// ack was received, close connection
	conn.Close()
}
