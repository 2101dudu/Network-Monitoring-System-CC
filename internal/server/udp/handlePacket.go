package udp

import (
	"fmt"
	"net"
	packet "nms/pkg/packet"
	utils "nms/pkg/utils"
)

func handlePacket(packetType utils.MessageType, packetPayload []byte, conn *net.UDPConn, udpAddr *net.UDPAddr) {
	switch packetType {
	case utils.ACK:
		fmt.Println("[SERVER] Acknowledgement received")
		return

	case utils.METRICSGATHERING:
		fmt.Println("[SERVER] Metrics received")
		return

	case utils.REGISTRATION:
		// TODO: change to thread
		fmt.Println("[SERVER] Processing registration request...")

		// Decode registration request
		reg, err := packet.DecodeRegistration(packetPayload)
		if err != nil {
			fmt.Println("[SERVER] [ERROR 12] Unable to decode registration data:", err)

			// send noack
			noack := packet.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).Build()
			packet.EncodeAndSendAck(conn, udpAddr, noack)
			return
		}

		// register agent
		mapOfAgents[reg.AgentID] = true

		// send ack
		ack := packet.NewAckBuilder().SetPacketID(reg.PacketID).SetSenderID(reg.AgentID).HasAcknowledged().Build()
		packet.EncodeAndSendAck(conn, udpAddr, ack)

		// Verify if isServer on any task - if so then sign all clients conected to that task that are already on the map!
		// If not server - verify if server already is running (after the agent server starts running iperf -s, sends to this server an ACK)
		// If server already running, send task to this client
		// If not, this client awaits with cond.wait (CHECK IT)
		// Maybe int that map of agents we can add more things like (just to make it easier):
		/* type Agent struct {
			AgentID byte
			IsServer bool
			Tasks []byte
		}

		var agents = make (map[byte] *Agent) with mutex */
		// If we the task just wants to know latency on this client, then we can just send it (because its just the command ping)

		return
	default:
		fmt.Println("[SERVER] [ERROR 13] Unknown message type")
		return
	}
}
