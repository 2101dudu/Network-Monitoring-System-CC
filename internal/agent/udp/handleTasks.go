package udp

import (
	"log"
	"net"
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

		if taskType != utils.PING && taskType != utils.IPERFCLIENT && taskType != utils.IPERFSERVER {
			log.Fatalln("[AGENT] [ERROR 80] Unexpected packet type received from server")
		}

		switch taskType {
		case utils.PING:
			handlePingTask(taskPayload)
		case utils.IPERFCLIENT:
			handleIperfClientTask(taskPayload, agentConn, udpAddr)
		case utils.IPERFSERVER:
			handleIperfServerTask(taskPayload, agentConn, udpAddr)
		}
	}
}

func handlePingTask(taskPayload []byte) {
	packet, err := task.DecodePingPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 81] Decoding ping packet")
	}

	// execute the pingPacket's command
	cmd := exec.Command("sh", "-c", packet.PingCommand)

	stdout, stderr := cmd.CombinedOutput()
	if stderr != nil {
		log.Fatalln("[AGENT] [ERROR 82] Executing ping command")
	}

	log.Println(string(stdout))
}

func handleIperfClientTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	iperfClient, err := task.DecodeIperfClientPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 81] Decoding ping packet")
	}

	// execute the pingPacket's command
	cmd := exec.Command(iperfClient.IperfClientCommand)

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

	// execute the pingPacket's command
	cmd := exec.Command(iperfServer.IperfServerCommand)

	data, err := cmd.Output()
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 82] Executing ping command")
	}

	log.Println(string(data))
}
