package udp

import (
	"log"
	"net"
	tcp "nms/internal/agent/tcp"
	ack "nms/internal/packet/ack"
	alert "nms/internal/packet/alert"
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
			handlePingTask(taskPayload, agentConn, udpAddr)
		case utils.IPERFCLIENT:
			handleIperfClientTask(taskPayload, agentConn, udpAddr)
		case utils.IPERFSERVER:
			handleIperfServerTask(taskPayload, agentConn, udpAddr)
		}
	}
}

func handlePingTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	packet, errDecode := task.DecodePingPacket(taskPayload)
	if errDecode != nil {
		log.Println("[AGENT] [ERROR 81] Decoding ping packet")
	}

	// IF VALIDAPACKET -> ENVIA ACK
	newAck := ack.NewAckBuilder().SetPacketID(packet.PacketID).SetSenderID(0).HasAcknowledged().Build()
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)
	// ELSE !VALIDAPACKET -> ENVIA NOACK, RETURN

	// execute the pingPacket's command
	output, cpuAlert, ramAlert, err := ExecuteCommandWithMonitoring(packet.PingCommand, packet.DeviceMetrics, packet.AlertFlowConditions)

	if err != nil {
		log.Println("[AGENT] [ERROR] Executing ping command")
	}

	if cpuAlert || ramAlert || err != nil {

		agentID, errAgent := utils.GetAgentID()
		if errAgent != nil {
			log.Fatalln("[AGENT] Unable to get agent ID:", errAgent)
		}

		buildAlert := alert.NewAlertBuilder().
			SetPacketID(packet.PacketID).
			SetSenderID(agentID).
			SetTaskID(packet.TaskID).
			SetCpuAlert(cpuAlert).
			SetRamAlert(ramAlert)

		if err != nil || errDecode != nil || errAgent != nil {
			buildAlert.SetErrorAlert(true)
		}
		//No iperf convém verificar no parse do output os outros dados como jitter e packetloss

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
	}

	//parse of output and send
	log.Println(output)
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
