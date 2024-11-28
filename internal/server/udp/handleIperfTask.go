package udp

import (
	"log"
	parse "nms/internal/jsonParse"
	t "nms/internal/packet/task"
	utils "nms/internal/utils"
)

func handleIperfTask(task parse.Task) {
	serverIndex, clientIndex := getIperfIndices(task)
	agentSIP := getAgentIP(task.Devices[serverIndex].DeviceID)
	agentCIP := getAgentIP(task.Devices[clientIndex].DeviceID)

	agentSConn := utils.ResolveUDPAddrAndDial(agentSIP, "9091")
	agentCConn := utils.ResolveUDPAddrAndDial(agentCIP, "9091")

	iperfServerPacket := t.ConvertTaskIntoIperfServerPacket(task, serverIndex)
	iperfClientPacket := t.ConvertTaskIntoIperfClientPacket(task, clientIndex)

	dataServer, err := t.EncodeIperfServerPacket(iperfServerPacket)
	if err != nil {
		log.Fatalln("[ERROR 25] Encoding iperf server packet")
	}

	dataClient, err := t.EncodeIperfClientPacket(iperfClientPacket)
	if err != nil {
		log.Fatalln("[ERROR 23] Encoding iperf client packet")
	}

	utils.WriteUDP(agentSConn, nil, dataServer, "[SERVER] [MAIN READ THREAD] Iperf server packet sent", "[SERVER] [ERROR 33] Unable to send iperf server packet")
	utils.WriteUDP(agentCConn, nil, dataClient, "[SERVER] [MAIN READ THREAD] Iperf client packet sent", "[SERVER] [ERROR 34] Unable to send iperf client packet")
}
