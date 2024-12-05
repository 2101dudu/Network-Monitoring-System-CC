package nettask

import (
	"log"
	"net"
	utils "nms/internal/utils"
	"sync"
)

var (
	metricsReceived = make(map[string]bool)
	metricsMutex    sync.Mutex
)

func handleMetrics(conn *net.UDPConn) {
	for {
		log.Println("[NetTask] Waiting for metrics from an agent")

		// Read metrics
		n, udpAddr, data := utils.ReadUDP(conn, "[NetTask] Metrics received", "[ERROR 10] Unable to read metrics")

		// Check if there is data
		if n == 0 {
			log.Println("[ERROR 11] No data received")
			continue
		}

		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		// Check if the packet type is correct
		if packetType != utils.METRICSGATHERING {
			log.Println("[ERROR 18] Unexpected packet type received from agent", packetType)
			continue
		}
		go handleMetricsGathering(packetPayload, conn, udpAddr)
	}
}
