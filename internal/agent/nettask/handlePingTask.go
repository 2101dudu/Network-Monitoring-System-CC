package nettask

import (
	"log"
	"net"
	tcp "nms/internal/agent/alertflow"
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
		log.Fatalln(utils.Red+"[ERROR 81] Decoding ping packet", utils.Reset)
	}

	agentID, errAgent := utils.GetAgentID()
	if errAgent != nil {
		log.Fatalln(utils.Red+"[ERROR 101] Unable to get agent ID:", errAgent, utils.Reset)
	}

	if !task.ValidateHashPingPacket(pingPacket) {
		noack := ack.NewAckBuilder().SetPacketID(pingPacket.PacketID).SetReceiverID(utils.SERVERID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(agentConn, udpAddr, noack)

		log.Println(utils.Red+"[ERROR 102] Invalid hash in ping packet", utils.Reset)
		return
	}

	newAck := ack.NewAckBuilder().SetPacketID(pingPacket.PacketID).SetReceiverID(utils.SERVERID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// Check if task was already received
	tasksMutex.Lock()
	if _, exists := tasksReceived[pingPacket.TaskID]; exists {
		tasksMutex.Unlock()
		return
	}
	tasksReceived[pingPacket.TaskID] = true
	tasksMutex.Unlock()

	availableTime := time.Duration(pingPacket.Frequency) * time.Second

	// reexecute the ping command every pingPacket.Frequency seconds
outerLoop:
	for {
		startTime := time.Now()

		for {
			// Check if the available time has passed
			if time.Since(startTime) > availableTime {
				// Send empty metrics packet and alert for timeout
				sendTimeoutAlertAndEmptyMetrics(pingPacket.PingCommand, pingPacket.TaskID, agentID, startTime)
				continue outerLoop
			}

			// Check if the iperf command is running
			iperfMutex.Lock()
			if !iperfRunning && time.Since(startTime) <= availableTime {
				iperfMutex.Unlock()
				break
			}
			iperfMutex.Unlock()

			time.Sleep(100 * time.Millisecond) // Wait a bit before checking again
		}

		// execute the pingPacket's command
		outputData, err := ExecuteCommandWithMonitoring(pingPacket.PingCommand, pingPacket.DeviceMetrics, pingPacket.AlertFlowConditions, pingPacket.TaskID)

		// Calculate the time that the command has left. This value can be negative if the command took longer than the frequency
		remainingIdleTime := availableTime - time.Since(startTime)

		// Send an alert if, during command execution, an error happened
		if err != nil {
			errTime := time.Now() // time of alert

			newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
			buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetSenderID(agentID).SetTaskID(pingPacket.TaskID).SetAlertType(alert.ERROR).SetTime(errTime.Format("15:04:05.000000000"))

			newAlert := buildAlert.Build()                        // build full alert with given sets
			tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
		}

		// If the remaining idle time is negative, send an empty metrics packet and alert for timeout
		if remainingIdleTime < 0 {
			sendTimeoutAlertAndEmptyMetrics(pingPacket.PingCommand, pingPacket.TaskID, agentID, startTime)
			continue
		}

		// Wait for the remaining idle time
		time.Sleep(remainingIdleTime)

		preparedOutput := parsePingOutput(string(outputData))

		serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

		metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(pingPacket.TaskID).SetTime(startTime.Format("15:04:05.000000000")).SetCommand(pingPacket.PingCommand).SetMetrics(preparedOutput).Build()

		hash = metrics.CreateHashMetricsPacket(newMetrics)
		newMetrics.Hash = (string(hash))

		packetData := metrics.EncodeMetrics(newMetrics)
		ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[NetTask] Metrics packet sent", "[ERROR 31] Unable to send metrics packet")
	}
}
