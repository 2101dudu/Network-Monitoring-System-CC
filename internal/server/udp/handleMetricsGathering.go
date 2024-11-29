package udp

import (
	"log"
	"net"
	"nms/internal/packet/ack"
	"nms/internal/packet/metrics"
)

func handleMetricsGathering(packetPayload []byte, conn *net.UDPConn, udpAddr *net.UDPAddr) {
	// Decode registration request
	reg, err := metrics.DecodeMetrics(packetPayload)
	if err != nil {
		log.Fatalln("[SERVER] [ERROR 12] Unable to decode metrics data:", err)
	}

	// TODO: CHECKSUM
	// noack := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
	// ack.EncodeAndSendAck(conn, udpAddr, noack)

	// send ack
	newAck := ack.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).HasAcknowledged().Build()
	ack.EncodeAndSendAck(conn, udpAddr, newAck)

	// store metrics
	//TODO
	log.Println("[SERVER] [METRICS] Metrics received from agent", reg.AgentID)
	log.Println("[SERVER] [METRICS] Metrics:", reg.Metrics)
}
