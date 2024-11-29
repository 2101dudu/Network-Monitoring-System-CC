package udp

import (
	"log"
	"net"
	ack "nms/internal/packet/ack"
	"nms/internal/packet/metrics"
	task "nms/internal/packet/task"
	utils "nms/internal/utils"
	"os/exec"
)

func handleTasks(agentConn *net.UDPConn) {
	for {
		n, udpAddr, taskData := utils.ReadUDP(agentConn, "[AGENT] [MAIN READ THREAD] Task received", "[AGENT] [ERROR 78] Unable to read task")
		if n == 0 {
			log.Println("[AGENT] [ERROR 79] No data received")
			continue
		}

		taskType := utils.PacketType(taskData[0])
		taskPayload := taskData[1:n]

		// Check if the packet type is correct
		if taskType != utils.PING && taskType != utils.IPERFCLIENT && taskType != utils.IPERFSERVER {
			log.Fatalln("[AGENT] [ERROR 80] Unexpected packet type received from server")
		}

		switch taskType {
		case utils.PING:
			handlePingTask(taskPayload, agentConn, udpAddr)
		case utils.IPERFCLIENT:
			handleIperfClientTask(taskPayload, agentConn, udpAddr)
		case utils.IPERFSERVER:
			handleIperfServerTask(taskPayload, agentConn, udpAddr)
		}
	}
}

func handlePingTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	pingPacket, err := task.DecodePingPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 81] Decoding ping packet")
	}

	// TODO: CHECKSUM
	// noack := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
	// ack.EncodeAndSendAck(conn, udpAddr, noack)

	newAck := ack.NewAckBuilder().SetPacketID(pingPacket.PacketID).SetSenderID(0).HasAcknowledged().Build()
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// execute the pingPacket's command
	cmd := exec.Command("sh", "-c", pingPacket.PingCommand)

	std, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 82] Executing ping command")
	}

	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

	var metricsID byte = 99
	newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetMetrics(string(std)).Build()
	data := metrics.EncodeMetrics(newMetrics)
	ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, data, "[SERVER] [MAIN READ THREAD] Metrics packet sent", "[SERVER] [ERROR 31] Unable to send metrics packet")
}

func handleIperfClientTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	iperfClient, err := task.DecodeIperfClientPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 81] Decoding ping packet")
	}

	// TODO: CHECKSUM
	// noack := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
	// ack.EncodeAndSendAck(conn, udpAddr, noack)

	newAck := ack.NewAckBuilder().SetPacketID(iperfClient.PacketID).SetSenderID(0).HasAcknowledged().Build()
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// execute the pingPacket's command
	cmd := exec.Command("sh", "-c", iperfClient.IperfClientCommand)

	data, err := cmd.Output()
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 82] Executing ping command")
	}

	log.Println(string(data))
}
func handleIperfServerTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	iperfServer, err := task.DecodeIperfServerPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 81] Decoding ping packet")
	}

	// TODO: CHECKSUM
	// noack := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
	// ack.EncodeAndSendAck(conn, udpAddr, noack)

	newAck := ack.NewAckBuilder().SetPacketID(iperfServer.PacketID).SetSenderID(0).HasAcknowledged().Build()
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// execute the pingPacket's command
	cmd := exec.Command("sh", "-c", iperfServer.IperfServerCommand)

	data, err := cmd.Output()
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 82] Executing ping command")
	}

	log.Println(string(data))
}
