package nettask

import (
	"log"
	"net"
	"nms/internal/utils"
	"sync"
)

var iperfRunning bool
var iperfMutex sync.Mutex

func handleTasks(agentConn *net.UDPConn) {
	for {
		n, udpAddr, taskData := utils.ReadUDP(agentConn, "[NetTask] Task received", "[ERROR 78] Unable to read task")
		if n == 0 {
			log.Println("[ERROR 79] No data received")
			continue
		}

		taskType := utils.PacketType(taskData[0])
		taskPayload := taskData[1:n]

		// Check if the packet type is correct
		if taskType != utils.PING && taskType != utils.IPERFCLIENT && taskType != utils.IPERFSERVER {
			log.Println("[ERROR 80] Unexpected packet type received from server")
			continue
		}

		switch taskType {
		case utils.PING:
			go handlePingTask(taskPayload, agentConn, udpAddr)
		case utils.IPERFCLIENT:
			go handleIperfClientTask(taskPayload, agentConn, udpAddr)
		case utils.IPERFSERVER:
			go handleIperfServerTask(taskPayload, agentConn, udpAddr)
		}
	}
}
