package nettask

import (
	"log"
	"net"
	tcp "nms/internal/agent/alertflow"
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
		log.Fatalln(utils.Red+"[ERROR 83] Decoding iperf server packet", utils.Reset)
	}

	if !task.ValidateHashIperfServerPacket(iperfServer) {
		noack := ack.NewAckBuilder().SetPacketID(iperfServer.PacketID).SetReceiverID(utils.SERVERID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(agentConn, udpAddr, noack)

		log.Println(utils.Red+"[ERROR 161] Invalid hash in iperf server packet", utils.Reset)
		return
	}
	newAck := ack.NewAckBuilder().SetPacketID(iperfServer.PacketID).SetReceiverID(utils.SERVERID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// Check if task was already received
	tasksMutex.Lock()
	if _, exists := tasksReceived[iperfServer.TaskID]; exists {
		tasksMutex.Unlock()
		return
	}
	tasksReceived[iperfServer.TaskID] = true
	tasksMutex.Unlock()

	availableTime := time.Duration(iperfServer.Frequency) * time.Second

	// reexecute the ping command every iperfServer.Frequency seconds
outerLoop:
	for {
		// keep track of the start time
		startTime := time.Now()

		for {
			// Check if the available time has passed
			if time.Since(startTime) > availableTime {
				// Send empty metrics packet and alert for timeout
				sendTimeoutAlertAndEmptyMetrics(iperfServer.IperfServerCommand, iperfServer.TaskID, agentID, startTime)
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
		outputData, err := ExecuteCommandWithMonitoring(iperfServer.IperfServerCommand, iperfServer.DeviceMetrics, iperfServer.AlertFlowConditions, iperfServer.TaskID)

		// Calculate the time that the command has left. This value can be negative if the command took longer than the frequency
		remainingIdleTime := availableTime - time.Since(startTime)

		errTime := time.Now() // time of alert

		// Send an alert if, during command execution, an error happened
		if err != nil {

			newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
			buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetAgentID(agentID).SetTaskID(iperfServer.TaskID).SetAlertType(alert.ERROR).SetTime(errTime.Format("15:04:05.000000000"))

			newAlert := buildAlert.Build()                        // build full alert with given sets
			tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
		}

		// Prepare output and check if jitter and packet loss exceeded
		preparedOutput, jitterHasExceeded, packetLossHasExceeded := parseIperfOutput(iperfServer.Bandwidth, iperfServer.Jitter, iperfServer.PacketLoss, float32(iperfServer.AlertFlowConditions.Jitter), float32(iperfServer.AlertFlowConditions.PacketLoss), string(outputData))

		if jitterHasExceeded > 1e-6 {
			newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
			buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetAgentID(agentID).SetTaskID(iperfServer.TaskID).SetAlertType(alert.JITTER).SetExceeded(jitterHasExceeded).SetTime(errTime.Format("15:04:05.000000000"))

			newAlert := buildAlert.Build()                        // build full alert with given sets
			tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
		}

		if packetLossHasExceeded > 1e-6 {
			newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
			buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetAgentID(agentID).SetTaskID(iperfServer.TaskID).SetAlertType(alert.PACKETLOSS).SetExceeded(packetLossHasExceeded).SetTime(errTime.Format("15:04:05.000000000"))

			newAlert := buildAlert.Build()                        // build full alert with given sets
			tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
		}

		// If the remaining idle time is negative, send an empty metrics packet and alert for timeout
		if remainingIdleTime < 0 {
			sendTimeoutAlertAndEmptyMetrics(iperfServer.IperfServerCommand, iperfServer.TaskID, agentID, startTime)
			continue
		}

		iperfMutex.Lock()
		iperfRunning = false
		iperfMutex.Unlock()

		// Wait for the remaining idle time
		time.Sleep(remainingIdleTime)

		serverConn := utils.ResolveUDPAddrAndDial(utils.SERVERIP, "8081")

		metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(iperfServer.TaskID).SetTime(startTime.Format("15:04:05.000000000")).SetCommand(iperfServer.IperfServerCommand).SetMetrics(preparedOutput).Build()

		hash = metrics.CreateHashMetricsPacket(newMetrics)
		newMetrics.Hash = (string(hash))

		packetData := metrics.EncodeMetrics(newMetrics)
		ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[NetTask] Metrics packet sent", "[ERROR 53] Unable to send metrics packet")
	}
}
