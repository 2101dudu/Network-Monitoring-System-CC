package nettask

import (
	"net"
	ack "nms/internal/packet/ack"
	registration "nms/internal/packet/registration"
	"nms/internal/utils"
)

func registerAgent(conn *net.UDPConn) {
	newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	registrationData := registration.CreateRegistrationPacket(newPacketID, agentID)

	successMessage := "[NetTask] Registration request sent"
	errorMessage := "[ERROR 159] Unable to send registration request"
	ack.SendPacketAndWaitForAck(newPacketID, agentID, packetsWaitingAck, &pMutex, conn, nil, registrationData, successMessage, errorMessage)
}
