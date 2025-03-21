package nettask

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"nms/internal/packet/ack"
	"nms/internal/packet/metrics"
	"nms/internal/utils"
	"os"
	"strconv"
	"sync"
)

var fileMutex sync.Mutex

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
		log.Fatalln(utils.Red+"[ERROR 153] Unable to decode metrics data:", err, utils.Reset)
	}

	if !metrics.ValidateHashMetricsPacket(met) {
		noack := ack.NewAckBuilder().SetPacketID(met.PacketID).SetReceiverID(met.AgentID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(conn, udpAddr, noack)

		log.Println(utils.Red+"[ERROR 100] Invalid hash in ping packet", utils.Reset)
		return
	}

	// send ack
	newAck := ack.NewAckBuilder().SetPacketID(met.PacketID).SetReceiverID(met.AgentID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(conn, udpAddr, newAck)

	// Check if metrics were already received
	mapID := fmt.Sprintf("%d:%d", met.PacketID, met.AgentID)

	metricsMutex.Lock()
	if _, exists := metricsReceived[mapID]; exists {
		return
	}
	metricsReceived[mapID] = true
	metricsMutex.Unlock()

	// store metrics
	metricsData := MetricsData{
		TaskID:       "task-" + strconv.Itoa(int(met.TaskID)),
		AgentID:      met.AgentID,
		LogTime:      met.Time,
		Command:      met.Command,
		OutputString: met.Metrics,
	}

	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile("output/metrics.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(utils.Red+"[ERROR 90] Unable to open metrics file:", err, utils.Reset)
	}
	defer file.Close()

	var metricsArray []MetricsData

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metricsArray); err != nil && err != io.EOF {
		log.Fatalln(utils.Red+"[ERROR 92] Unable to decode metrics data:", err, utils.Reset)
	}

	metricsArray = append(metricsArray, metricsData)

	file.Seek(0, 0)
	file.Truncate(0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Set indentation for pretty-printing
	if err := encoder.Encode(metricsArray); err != nil {
		log.Fatalln(utils.Red+"[ERROR 91] Unable to encode metrics data:", err, utils.Reset)
	}
}
