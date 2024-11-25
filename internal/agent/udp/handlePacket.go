package udp

import (
	"log"
	"net"
	ack "nms/internal/packet/ack"
	utils "nms/internal/utils"
)

func handlePacket(packetType utils.PacketType, packetPayload []byte, conn *net.UDPConn) {
	switch packetType {
	case utils.ACK:
		ack.HandleAck(packetPayload, packetsWaitingAck, &pMutex, agentID, conn)
		return
	case utils.PING:
		log.Println("[AGENT] Metrics received from server")
		// HandleTask method - TO DO
		return
	default:
		log.Println("[AGENT] [ERROR 7] Unknown packet type received from server")
		return
	}
}
