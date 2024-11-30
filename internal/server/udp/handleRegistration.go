package udp

import (
	"log"
	"net"
	ack "nms/internal/packet/ack"
	registration "nms/internal/packet/registration"
)

func handleRegistration(packetPayload []byte, conn *net.UDPConn, udpAddr *net.UDPAddr) {
	// Decode registration request
	reg, err := registration.DecodeRegistration(packetPayload)
	if err != nil {
		log.Fatalln("[SERVER] [ERROR 12] Unable to decode registration data:", err)
	}

	if !registration.ValidateHashRegistrationPacket(reg) {
		noack := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(conn, udpAddr, noack)

		log.Println("[SERVER] [ERROR 99] Invalid hash in registration packet")
		return
	}

	// register agent
	agentsIPs[reg.AgentID] = reg.IP

	// send ack
	newAck := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(conn, udpAddr, newAck)
}
