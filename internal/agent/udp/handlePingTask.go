package udp

import (
	"log"
	"net"
	ack "nms/internal/packet/ack"
	"nms/internal/packet/metrics"
	"nms/internal/packet/task"
	"nms/internal/utils"
	"os/exec"
)

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
