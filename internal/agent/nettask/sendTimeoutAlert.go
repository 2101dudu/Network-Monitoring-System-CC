package nettask

import (
	tcp "nms/internal/agent/alertflow"
	ack "nms/internal/packet/ack"
	alert "nms/internal/packet/alert"
	metrics "nms/internal/packet/metrics"
	utils "nms/internal/utils"
	"time"
)

func sendTimeoutAlertAndEmptyMetrics(command string, taskID uint16, agentID byte, startTime time.Time) {
	// Send alert for timeout
	newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetSenderID(agentID).SetTaskID(taskID).SetAlertType(alert.TIMEOUT).SetTime(startTime.Format("15:04:05.000000000"))

	newAlert := buildAlert.Build()
	tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert)

	// Send empty metrics packet
	serverConn := utils.ResolveUDPAddrAndDial(utils.SERVERIP, "8081")
	metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(taskID).SetTime(startTime.Format("15:04:05.000000000")).SetCommand(command).SetMetrics("").Build()

	hash := metrics.CreateHashMetricsPacket(newMetrics)
	newMetrics.Hash = (string(hash))

	packetData := metrics.EncodeMetrics(newMetrics)
	ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[NetTask] Empty metrics packet sent", "[ERROR 331] Unable to send empty metrics packet")
}
