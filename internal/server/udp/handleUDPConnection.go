package udp

import (
	"fmt"
	"net"
	utils "nms/pkg/utils"
)

var agentsIPs map[byte][4]byte

func handleUDPConnection(conn *net.UDPConn) {
	fmt.Println("[SERVER] [MAIN READ THREAD] Waiting for data from an agent")

	// Read registration request
	n, udpAddr, data := utils.ReadUDP(conn, "[SERVER] [MAIN READ THREAD] Registration request received", "[SERVER] [ERROR 10] Unable to read registration request")

	// Check if there is data
	if n == 0 {
		fmt.Println("[SERVER] [MAIN READ THREAD] [ERROR 11] No data received")
		return
	}

	// type cast the data to the appropriate message type
	packetType := utils.MessageType(data[0])
	packetPayload := data[1:n]

	handlePacket(packetType, packetPayload, conn, udpAddr)
}
