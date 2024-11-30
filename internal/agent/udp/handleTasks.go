package udp

import (
	"log"
	"net"
	"nms/internal/utils"
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
