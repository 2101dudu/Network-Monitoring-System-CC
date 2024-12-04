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

func handleIperfClientTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	iperfClient, err := task.DecodeIperfClientPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 85] Decoding iperf client packet")
	}

	if !task.ValidateHashIperfClientPacket(iperfClient) {
		noack := ack.NewAckBuilder().SetPacketID(iperfClient.PacketID).SetReceiverID(utils.SERVERID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(agentConn, udpAddr, noack)

		log.Println("[AGENT] [ERROR 100] Invalid hash in iperf client packet")
		return
	}

	newAck := ack.NewAckBuilder().SetPacketID(iperfClient.PacketID).SetReceiverID(utils.SERVERID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	availableTime := time.Duration(iperfClient.Frequency) * time.Second

	// reexecute the ping command every iperfClient.Frequency seconds
outerLoop:
	for {
		startTime := time.Now()

		for {
			// Check if the available time has passed
			if time.Since(startTime) > availableTime {
				// Send empty metrics packet and alert for timeout
				sendTimeoutAlertAndEmptyMetrics(iperfClient.IperfClientCommand, iperfClient.TaskID, agentID, startTime)
				continue outerLoop
			}

			// Check if another iperf command is running
			iperfMutex.Lock()
			if !iperfRunning && time.Since(startTime) <= availableTime {
				iperfRunning = true // Set iperfRunning to true to prevent other iperf commands from running
				iperfMutex.Unlock()
				break
			}
			iperfMutex.Unlock()

			time.Sleep(100 * time.Millisecond) // Wait a bit before checking again
		}

		// execute the iperfPacket's command
		outputData, err := ExecuteCommandWithMonitoring(iperfClient.IperfClientCommand, iperfClient.DeviceMetrics, iperfClient.AlertFlowConditions, iperfClient.TaskID)

		// Calculate the time that the command has left. This value can be negative if the command took longer than the frequency
		remainingIdleTime := availableTime - time.Since(startTime)

		errTime := time.Now() // time of alert

		// Send an alert if, during command execution, an error happened
		if err != nil {
			newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
			buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetSenderID(agentID).SetTaskID(iperfClient.TaskID).SetAlertType(alert.ERROR).SetTime(errTime.Format("15:04:05.000000000"))

			newAlert := buildAlert.Build()                        // build full alert with given sets
			tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
		}

		// Prepare output and check if jitter and packet loss exceeded
		preparedOutput, jitterHasExceeded, packetLossHasExceeded := parseIperfOutput(iperfClient.Bandwidth, iperfClient.Jitter, iperfClient.PacketLoss, float32(iperfClient.AlertFlowConditions.Jitter), float32(iperfClient.AlertFlowConditions.PacketLoss), string(outputData))

		if jitterHasExceeded > 1e-6 {
			newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
			buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetSenderID(agentID).SetTaskID(iperfClient.TaskID).SetAlertType(alert.JITTER).SetExceeded(jitterHasExceeded).SetTime(errTime.Format("15:04:05.000000000"))

			newAlert := buildAlert.Build()                        // build full alert with given sets
			tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
		}

		if packetLossHasExceeded > 1e-6 {
			newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
			buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetSenderID(agentID).SetTaskID(iperfClient.TaskID).SetAlertType(alert.PACKETLOSS).SetExceeded(packetLossHasExceeded).SetTime(errTime.Format("15:04:05.000000000"))

			newAlert := buildAlert.Build()                        // build full alert with given sets
			tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
		}

		// If the remaining idle time is negative, send an empty metrics packet and alert for timeout
		if remainingIdleTime < 0 {
			sendTimeoutAlertAndEmptyMetrics(iperfClient.IperfClientCommand, iperfClient.TaskID, agentID, startTime)
			continue
		}

		iperfMutex.Lock()
		iperfRunning = false
		iperfMutex.Unlock()

		// Wait for the remaining idle time
		time.Sleep(remainingIdleTime)

		serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

		metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(iperfClient.TaskID).SetTime(startTime.Format("15:04:05.000000000")).SetCommand(iperfClient.IperfClientCommand).SetMetrics(preparedOutput).Build()

		hash = metrics.CreateHashMetricsPacket(newMetrics)
		newMetrics.Hash = (string(hash))

		packetData := metrics.EncodeMetrics(newMetrics)
		ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[SERVER] [MAIN READ THREAD] Metrics packet sent", "[SERVER] [ERROR 36] Unable to send metrics packet")
	}
}
