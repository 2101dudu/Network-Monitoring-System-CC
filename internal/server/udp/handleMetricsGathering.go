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
)

// MetricsData represents the structure of the metrics data to be stored in the JSON file
type MetricsData struct {
	TaskID       string `json:"task_id"`
	AgentID      byte   `json:"agent_id"`
	LogTime      string `json:"log_time"`
	Command      string `json:"command"`
	OutputString string `json:"output_string"`
}

func handleMetricsGathering(packetPayload []byte, conn *net.UDPConn, udpAddr *net.UDPAddr) {
	// Decode registration request
	met, err := metrics.DecodeMetrics(packetPayload)
	if err != nil {
		log.Fatalln("[SERVER] [ERROR 12] Unable to decode metrics data:", err)
	}

	if !metrics.ValidateHashMetricsPacket(met) {
		noack := ack.NewAckBuilder().SetPacketID(met.PacketID).SetReceiverID(met.AgentID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(conn, udpAddr, noack)

		log.Println("[AGENT] [ERROR 100] Invalid hash in ping packet")
		return
	}

	// send ack
	newAck := ack.NewAckBuilder().SetPacketID(met.PacketID).SetReceiverID(met.AgentID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(conn, udpAddr, newAck)

	// store metrics
	metricsData := MetricsData{
		TaskID:       "task-" + strconv.Itoa(int(met.TaskID)),
		AgentID:      met.AgentID,
		LogTime:      met.Time,
		Command:      met.Command,
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
