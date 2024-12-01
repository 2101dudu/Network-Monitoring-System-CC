package udp

import (
	"log"
	"net"
	ack "nms/internal/packet/ack"
	"nms/internal/packet/metrics"
	"nms/internal/packet/task"
	"nms/internal/utils"
	"os/exec"
	"time"
)

func handlePingTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	pingPacket, err := task.DecodePingPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 81] Decoding ping packet")
	}

	if !task.ValidateHashPingPacket(pingPacket) {
		noack := ack.NewAckBuilder().SetPacketID(pingPacket.PacketID).SetSenderID(utils.SERVERID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(agentConn, udpAddr, noack)

		log.Println("[AGENT] [ERROR 102] Invalid hash in ping packet")
		return
	}

	newAck := ack.NewAckBuilder().SetPacketID(pingPacket.PacketID).SetSenderID(utils.SERVERID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// keep track of the start time
	startTime := time.Now()

	// execute the pingPacket's command
	cmd := exec.Command("sh", "-c", pingPacket.PingCommand)

	outputData, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 82] Executing ping command", err)
	}

	preparedOutput := parsePingOutput(string(outputData))

	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

	metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(pingPacket.TaskID).SetTime(startTime.Format("15:04:05.000000000")).SetMetrics(preparedOutput).Build()

	hash = metrics.CreateHashMetricsPacket(newMetrics)
	newMetrics.Hash = (string(hash))

	packetData := metrics.EncodeMetrics(newMetrics)
	ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[SERVER] [MAIN READ THREAD] Metrics packet sent", "[SERVER] [ERROR 31] Unable to send metrics packet")
}
