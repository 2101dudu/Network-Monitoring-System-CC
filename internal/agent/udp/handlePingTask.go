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
	outputData, cpuAlert, ramAlert, err := ExecuteCommandWithMonitoring(pingPacket.PingCommand, pingPacket.DeviceMetrics, pingPacket.AlertFlowConditions)

	if err != nil {
		log.Println("[AGENT] [ERROR] Executing ping command")
	}

	if cpuAlert || ramAlert || err != nil {

		agentID, errAgent := utils.GetAgentID()
		if errAgent != nil {
			log.Fatalln("[AGENT] Unable to get agent ID:", errAgent)
		}

		buildAlert := alert.NewAlertBuilder().
			SetPacketID(pingPacket.PacketID).
			SetSenderID(agentID).
			SetTaskID(pingPacket.TaskID).
			SetCpuAlert(cpuAlert).
			SetRamAlert(ramAlert)

		if err != nil || errAgent != nil {
			buildAlert.SetErrorAlert(true)
		}
		//No iperf convém verificar no parse do output os outros dados como jitter e packetloss

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
	}

	preparedOutput := parsePingOutput(string(outputData))

	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

	metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTime(startTime.Format("15:04:05.000000000")).SetMetrics(preparedOutput).Build()

	hash = metrics.CreateHashMetricsPacket(newMetrics)
	newMetrics.Hash = (string(hash))

	packetData := metrics.EncodeMetrics(newMetrics)
	ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[SERVER] [MAIN READ THREAD] Metrics packet sent", "[SERVER] [ERROR 31] Unable to send metrics packet")
}
