package udp

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"nms/internal/packet/ack"
	"nms/internal/packet/metrics"
	"os"
	"strconv"
	"time"
)

// MetricsData represents the structure of the metrics data to be stored in the JSON file
type MetricsData struct {
	TaskID       string `json:"task_id"`
	AgentID      byte   `json:"agent_id"`
	LogTime      string `json:"log_time"`
	OutputString string `json:"output_string"`
}

func handleMetricsGathering(packetPayload []byte, conn *net.UDPConn, udpAddr *net.UDPAddr) {
	// Decode registration request
	met, err := metrics.DecodeMetrics(packetPayload)
	if err != nil {
		log.Fatalln("[SERVER] [ERROR 12] Unable to decode metrics data:", err)
	}

	// TODO: CHECKSUM
	// noack := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
	// ack.EncodeAndSendAck(conn, udpAddr, noack)

	// send ack
	newAck := ack.NewAckBuilder().SetPacketID(met.PacketID).SetSenderID(met.AgentID).HasAcknowledged().Build()
	ack.EncodeAndSendAck(conn, udpAddr, newAck)

	// store metrics
	metricsData := MetricsData{
		TaskID:       "task-" + strconv.Itoa(int(met.TaskID)),
		AgentID:      met.AgentID,
		LogTime:      time.Now().Format("15:04:05.000000000"),
		OutputString: met.Metrics,
	}

	file, err := os.OpenFile("output/metrics.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("[SERVER] [ERROR 90] Unable to open metrics file:", err)
	}
	defer file.Close()

	var metricsArray []MetricsData

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metricsArray); err != nil && err != io.EOF {
		log.Fatalln("[SERVER] [ERROR 92] Unable to decode metrics data:", err)
	}

	metricsArray = append(metricsArray, metricsData)

	file.Seek(0, 0)
	file.Truncate(0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Set indentation for pretty-printing
	if err := encoder.Encode(metricsArray); err != nil {
		log.Fatalln("[SERVER] [ERROR 91] Unable to encode metrics data:", err)
	}
}
