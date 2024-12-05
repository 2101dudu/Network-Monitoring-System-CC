package nettask

import (
	"log"
	"net"
	ack "nms/internal/packet/ack"
	registration "nms/internal/packet/registration"
	"nms/internal/utils"
)

func handleRegistration(packetPayload []byte, conn *net.UDPConn, udpAddr *net.UDPAddr) {
	// Decode registration request
	reg, err := registration.DecodeRegistration(packetPayload)
	if err != nil {
		log.Fatalln(utils.Red+"[ERROR 12] Unable to decode registration data:", err, utils.Reset)
	}

	if !registration.ValidateHashRegistrationPacket(reg) {
		noack := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetReceiverID(reg.AgentID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(conn, udpAddr, noack)

		log.Println(utils.Red+"[ERROR 99] Invalid hash in registration packet", utils.Reset)
		return
	}

	// register agent
	agentsIPs[reg.AgentID] = reg.IP

	// send ack
	newAck := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetReceiverID(reg.AgentID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(conn, udpAddr, newAck)
}
