package nettask

import (
	"log"
	"net"
	utils "nms/internal/utils"
)

func handleRegistrations(conn *net.UDPConn) {
	log.Println(utils.Blue, "Waiting for", numAgents, "agent(s) to register", utils.Reset)

	for len(agentsIPs) < numAgents {
		log.Println(utils.Blue, "Total agents registered until now:", len(agentsIPs), utils.Reset)

		// Read registration request
		n, udpAddr, data := utils.ReadUDP(conn, "[ERROR 10] Unable to read registration request")

		// Check if there is data
		if n == 0 {
			log.Println(utils.Red, "[ERROR 11] No data received", utils.Reset)
			continue
		}

		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		// Check if the packet type is correct
		if packetType != utils.REGISTRATION {
			log.Println(utils.Red, "[ERROR 18] Unexpected packet type received from server", utils.Reset)
			continue
		}
		handleRegistration(packetPayload, conn, udpAddr)
	}
}
