package udp

import (
	"fmt"
	"net"
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
		fmt.Println("[SERVER] Processing registration request...")
		handleRegistration(packetPayload, conn, udpAddr)
		return
	default:
		fmt.Println("[SERVER] [ERROR 13] Unknown message type")
		return
	}
}
