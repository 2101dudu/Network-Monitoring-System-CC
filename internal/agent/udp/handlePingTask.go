package udp

import (
	"log"
	"net"
	tcp "nms/internal/agent/tcp"
	ack "nms/internal/packet/ack"
	alert "nms/internal/packet/alert"
	metrics "nms/internal/packet/metrics"
	task "nms/internal/packet/task"
	utils "nms/internal/utils"
	"time"
)

func handlePingTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	pingPacket, err := task.DecodePingPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 81] Decoding ping packet")
	}

	agentID, errAgent := utils.GetAgentID()
	if errAgent != nil {
		log.Fatalln("[AGENT] [ERROR 101] Unable to get agent ID:", errAgent)
	}

	if !task.ValidateHashPingPacket(pingPacket) {
		noack := ack.NewAckBuilder().SetPacketID(pingPacket.PacketID).SetReceiverID(utils.SERVERID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(agentConn, udpAddr, noack)

		log.Println("[AGENT] [ERROR 102] Invalid hash in ping packet")
		return
	}

	newAck := ack.NewAckBuilder().SetPacketID(pingPacket.PacketID).SetReceiverID(utils.SERVERID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// keep track of the start time
	startTime := time.Now()

	// execute the pingPacket's command
	outputData, err := ExecuteCommandWithMonitoring(pingPacket.PingCommand, pingPacket.DeviceMetrics, pingPacket.AlertFlowConditions, pingPacket.TaskID)

	if err != nil { // If during command execution happened an error then send an alert
		errTime := time.Now() // time of alert

		newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		buildAlert := alert.NewAlertBuilder().
			SetPacketID(newPacketID).
			SetSenderID(agentID).
			SetTaskID(pingPacket.TaskID).
			SetAlertType(alert.ERROR).
			SetTime(errTime.Format("15:04:05.000000000"))

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
	}

	preparedOutput := parsePingOutput(string(outputData))

	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

	metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(pingPacket.TaskID).SetTime(startTime.Format("15:04:05.000000000")).SetCommand(pingPacket.PingCommand).SetMetrics(preparedOutput).Build()

	hash = metrics.CreateHashMetricsPacket(newMetrics)
	newMetrics.Hash = (string(hash))

	packetData := metrics.EncodeMetrics(newMetrics)
	ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[SERVER] [MAIN READ THREAD] Metrics packet sent", "[SERVER] [ERROR 31] Unable to send metrics packet")
}
