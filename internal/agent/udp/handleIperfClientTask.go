package udp

import (
	"log"
	"net"
	ack "nms/internal/packet/ack"
	"nms/internal/packet/task"
	"os/exec"
)

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
