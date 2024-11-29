package udp

import (
	"net"
	ack "nms/internal/packet/ack"
	registration "nms/internal/packet/registration"
)

var agentID byte

func registerAgent(conn *net.UDPConn, agentIP string) {
	var firstPacketID byte = 1
	var registrationData []byte
	agentID, registrationData = registration.CreateRegistrationPacket(firstPacketID, agentIP)

	successMessage := "[AGENT] Registration request sent"
	errorMessage := "[AGENT] [ERROR 4] Unable to send registration request"
	ack.SendPacketAndWaitForAck(firstPacketID, agentID, packetsWaitingAck, &pMutex, conn, nil, registrationData, successMessage, errorMessage)
}
