package udp

import (
	"log"
	"net"
	utils "nms/internal/utils"
)

func handleMetrics(conn *net.UDPConn) {
	for {
		log.Println("[SERVER] [MAIN READ THREAD] Waiting for metrics from an agent")

		// Read metrics
		n, udpAddr, data := utils.ReadUDP(conn, "[SERVER] [MAIN READ THREAD] Metrics received", "[SERVER] [ERROR 10] Unable to read metrics")

		// Check if there is data
		if n == 0 {
			log.Println("[SERVER] [MAIN READ THREAD] [ERROR 11] No data received")
			continue
		}

		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		// Check if the packet type is correct
		if packetType != utils.METRICSGATHERING {
			log.Println(packetType)
			log.Fatalln("[AGENT] [ERROR 18] Unexpected packet type received from server")
		}
		handleMetricsGathering(packetPayload, conn, udpAddr)
	}
}
