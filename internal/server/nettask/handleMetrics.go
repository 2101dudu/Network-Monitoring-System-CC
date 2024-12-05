package nettask

import (
	"log"
	"net"
	utils "nms/internal/utils"
)

func handleMetrics(conn *net.UDPConn) {
	for {

		// Read metrics
		n, udpAddr, data := utils.ReadUDP(conn, "[ERROR 10] Unable to read metrics")

		// Check if there is data
		if n == 0 {
			log.Println(utils.Red, "[ERROR 11] No data received", utils.Reset)
			continue
		}

		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		// Check if the packet type is correct
		if packetType != utils.METRICSGATHERING {
			log.Println(utils.Red, "[ERROR 18] Unexpected packet type received from agent", packetType, utils.Reset)
			continue
		}
		handleMetricsGathering(packetPayload, conn, udpAddr)
		log.Println(utils.Blue, "[NetTask] Metrics received", utils.Reset)
	}
}
