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
		log.Println(utils.Blue+"[NetTask] Waiting for metrics from an agent", utils.Reset)

		// Read metrics
		n, udpAddr, data := utils.ReadUDP(conn, "[NetTask] Metrics received", "[ERROR 61] Unable to read metrics")

		// Check if there is data
		if n == 0 {
			log.Println(utils.Red+"[ERROR 11] No data received", utils.Reset)
			continue
		}

		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		// Check if the packet type is correct
		if packetType != utils.METRICSGATHERING {
			log.Println(utils.Red+"[ERROR 152] Unexpected packet type received from agent", packetType, utils.Reset)
			continue
		}
		go handleMetricsGathering(packetPayload, conn, udpAddr)
	}
}
