package udp

import (
	"fmt"
	"net"
	packet "nms/internal/packet"
	utils "nms/pkg/utils"
)

func handlePacket(packetType utils.PacketType, packetPayload []byte, conn *net.UDPConn) {
	switch packetType {
	case utils.ACK:
		packet.HandleAck(packetPayload, packetsWaitingAck, &pMutex, agentID, conn)
		return
	case utils.TASK:
		fmt.Println("[AGENT] Metrics received from server")
		// HandleTask method - TO DO
		return
	default:
		fmt.Println("[AGENT] [ERROR 7] Unknown message type received from server")
		return
	}
}
