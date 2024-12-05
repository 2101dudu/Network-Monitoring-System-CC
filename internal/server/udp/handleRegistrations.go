package udp

import (
	"fmt"
	"log"
	"net"
	utils "nms/internal/utils"
)

func handleRegistrations(conn *net.UDPConn) {
	fmt.Printf("Waiting for %d agents to register\n", numAgents)

	for len(agentsIPs) < numAgents {
		fmt.Printf("Total agents registered until now: %d\n", len(agentsIPs))

		// Read registration request
		n, udpAddr, data := utils.ReadUDP(conn, "[NetTask] Registration request received", "[ERROR 10] Unable to read registration request")

		// Check if there is data
		if n == 0 {
			log.Println("[ERROR 11] No data received")
			continue
		}

		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		// Check if the packet type is correct
		if packetType != utils.REGISTRATION {
			log.Println("[ERROR 18] Unexpected packet type received from server")
			continue
		}
		handleRegistration(packetPayload, conn, udpAddr)
	}
}
