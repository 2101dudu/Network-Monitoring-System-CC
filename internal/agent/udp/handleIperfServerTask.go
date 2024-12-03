package udp

import (
	"log"
	"net"
	tcp "nms/internal/agent/tcp"
	ack "nms/internal/packet/ack"
	alert "nms/internal/packet/alert"
	"nms/internal/packet/metrics"
	"nms/internal/packet/task"
	"nms/internal/utils"
	"time"
)

func handleIperfServerTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	iperfServer, err := task.DecodeIperfServerPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 83] Decoding iperf server packet")
	}

	if !task.ValidateHashIperfServerPacket(iperfServer) {
		noack := ack.NewAckBuilder().SetPacketID(iperfServer.PacketID).SetReceiverID(utils.SERVERID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(agentConn, udpAddr, noack)

		log.Println("[AGENT] [ERROR 101] Invalid hash in iperf server packet")
		return
	}
	newAck := ack.NewAckBuilder().SetPacketID(iperfServer.PacketID).SetReceiverID(utils.SERVERID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// reexecute the ping command every iperfServer.Frequency seconds
	for {
		// keep track of the start time
		startTime := time.Now()

	// execute the iperfPacket's command
	outputData, err := ExecuteCommandWithMonitoring(iperfServer.IperfServerCommand, iperfServer.DeviceMetrics, iperfServer.AlertFlowConditions, iperfServer.TaskID)

	if err != nil { // If during command execution happened an error, send an alert

		newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		buildAlert := alert.NewAlertBuilder().
			SetPacketID(newPacketID).
			SetSenderID(agentID).
			SetTaskID(iperfServer.TaskID).
			SetAlertType(alert.ERROR).
			SetTime(startTime.Format("15:04:05.000000000"))

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
	}

	// Prepare output and check if jitter and packet loss exceeded
	preparedOutput, jitterHasExceeded, packetLossHasExceeded := parseIperfOutput(iperfServer.Bandwidth, iperfServer.Jitter, iperfServer.PacketLoss, float32(iperfServer.AlertFlowConditions.Jitter), float32(iperfServer.AlertFlowConditions.PacketLoss), string(outputData))

	if jitterHasExceeded > 1e-6 {

		newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		buildAlert := alert.NewAlertBuilder().
			SetPacketID(newPacketID).
			SetSenderID(agentID).
			SetTaskID(iperfServer.TaskID).
			SetAlertType(alert.JITTER).
			SetExceeded(jitterHasExceeded).
			SetTime(startTime.Format("15:04:05.000000000"))

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
	}

	if packetLossHasExceeded > 1e-6 {
		newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		buildAlert := alert.NewAlertBuilder().
			SetPacketID(newPacketID).
			SetSenderID(agentID).
			SetTaskID(iperfServer.TaskID).
			SetAlertType(alert.PACKETLOSS).
			SetExceeded(packetLossHasExceeded).
			SetTime(startTime.Format("15:04:05.000000000"))

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
	}

		// calculate the elapsed time and sleep for the remaining time to ensure the loop runs every iperfServer.Frequency seconds
		elapsedTime := time.Since(startTime)
		sleepDuration := time.Duration(iperfServer.Frequency)*time.Second - elapsedTime
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}

		serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")
    
	  metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	  newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(iperfServer.TaskID).SetTime(startTime.Format("15:04:05.000000000")).SetCommand(iperfServer.IperfServerCommand).SetMetrics(preparedOutput).Build()
    
		hash = metrics.CreateHashMetricsPacket(newMetrics)
		newMetrics.Hash = (string(hash))

		packetData := metrics.EncodeMetrics(newMetrics)
		ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[SERVER] [MAIN READ THREAD] Metrics packet sent", "[SERVER] [ERROR 35] Unable to send metrics packet")
	}
}
