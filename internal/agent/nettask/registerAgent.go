package nettask

import (
	"net"
	ack "nms/internal/packet/ack"
	registration "nms/internal/packet/registration"
	"nms/internal/utils"
)

var agentID byte

func registerAgent(conn *net.UDPConn, agentIP string) {
	newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	var registrationData []byte
	agentID, registrationData = registration.CreateRegistrationPacket(newPacketID, agentIP)

	successMessage := "[NetTask] Registration request sent"
	errorMessage := "[ERROR 159] Unable to send registration request"
	ack.SendPacketAndWaitForAck(newPacketID, agentID, packetsWaitingAck, &pMutex, conn, nil, registrationData, successMessage, errorMessage)
}
