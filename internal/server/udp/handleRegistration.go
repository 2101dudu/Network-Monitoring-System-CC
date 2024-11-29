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

	// TODO: CHECKSUM
	// noack := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
	// ack.EncodeAndSendAck(conn, udpAddr, noack)

	// register agent
	agentsIPs[reg.AgentID] = reg.IP

	// send ack
	newAck := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).HasAcknowledged().Build()
	ack.EncodeAndSendAck(conn, udpAddr, newAck)
}
