package udp

import (
	"fmt"
	"net"
	packet "nms/pkg/packet"
	utils "nms/pkg/utils"
)

func handlePacket(packetType utils.MessageType, packetPayload []byte, conn *net.UDPConn, udpAddr *net.UDPAddr) {
	switch packetType {
	case utils.ACK:
		fmt.Println("[SERVER] Acknowledgement received")
		return

	case utils.METRICSGATHERING:
		fmt.Println("[SERVER] Metrics received")
		return

	case utils.REGISTRATION:
		// TODO: change to thread
		fmt.Println("[SERVER] Processing registration request...")

		// Decode registration request
		reg, err := packet.DecodeRegistration(packetPayload)
		if err != nil {
			fmt.Println("[SERVER] [ERROR 12] Unable to decode registration data:", err)

			// send noack
			noack := packet.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
			packet.EncodeAndSendAck(conn, udpAddr, noack)
			return
		}

		// register agent
		mapOfAgents[reg.AgentID] = true

		// send ack
		ack := packet.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).HasAcknowledged().Build()
		packet.EncodeAndSendAck(conn, udpAddr, ack)
		return
	default:
		fmt.Println("[SERVER] [ERROR 13] Unknown message type")
		return
	}
}
