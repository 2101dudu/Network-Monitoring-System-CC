package udp

import (
	"fmt"
	"net"
	utils "nms/internal/utils"
)

func handleRegistrations(conn *net.UDPConn) {
	for len(agentsIPs) < utils.MAXAGENTS {
		fmt.Println("[SERVER] [MAIN READ THREAD] Waiting for data from an agent")

		// Read registration request
		n, udpAddr, data := utils.ReadUDP(conn, "[SERVER] [MAIN READ THREAD] Registration request received", "[SERVER] [ERROR 10] Unable to read registration request")

		// Check if there is data
		if n == 0 {
			fmt.Println("[SERVER] [MAIN READ THREAD] [ERROR 11] No data received")
			return
		}

		// type cast the data to the appropriate message type
		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		if packetType != utils.REGISTRATION {
			fmt.Println("[AGENT] [ERROR 18] Unexpected message type received from server")
			return
		}
		handleRegistration(packetPayload, conn, udpAddr)
	}
	conn.Close()
}
