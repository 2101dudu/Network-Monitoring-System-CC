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

	// reexecute the ping command every iperfClient.Frequency seconds
	for {
		// keep track of the start time
		startTime := time.Now()

		// execute the iperf client packets's command
		cmd := exec.Command("sh", "-c", iperfClient.IperfClientCommand)

		outputData, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalln("[AGENT] [ERROR 86] Executing iperf client command")
		}

		preparedOutput := parseIperfOutput(iperfClient.Bandwidth, iperfClient.Jitter, iperfClient.PacketLoss, string(outputData))

		// calculate the elapsed time and sleep for the remaining time to ensure the loop runs every iperfClient.Frequency seconds
		elapsedTime := time.Since(startTime)
		sleepDuration := time.Duration(iperfClient.Frequency)*time.Second - elapsedTime
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}

		serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

		metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(iperfClient.TaskID).SetTime(startTime.Format("15:04:05.000000000")).SetMetrics(preparedOutput).Build()

		hash = metrics.CreateHashMetricsPacket(newMetrics)
		newMetrics.Hash = (string(hash))

		packetData := metrics.EncodeMetrics(newMetrics)
		ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[SERVER] [MAIN READ THREAD] Metrics packet sent", "[SERVER] [ERROR 36] Unable to send metrics packet")
	}
}
