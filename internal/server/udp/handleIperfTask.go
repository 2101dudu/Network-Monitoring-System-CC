package udp

import (
	"log"
	parse "nms/internal/jsonParse"
	"nms/internal/packet/ack"
	t "nms/internal/packet/task"
	utils "nms/internal/utils"
)

func handleIperfTask(task parse.Task) {
	serverIndex, clientIndex := getIperfIndices(task)
	agentSIP := getAgentIP(task.Devices[serverIndex].DeviceID)
	agentCIP := getAgentIP(task.Devices[clientIndex].DeviceID)

	agentSConn := utils.ResolveUDPAddrAndDial(agentSIP, "9091")
	agentCConn := utils.ResolveUDPAddrAndDial(agentCIP, "9091")

	iperfServerPacket := ConvertTaskIntoIperfServerPacket(task, serverIndex)
	iperfClientPacket := ConvertTaskIntoIperfClientPacket(task, clientIndex)

	dataServer, err := t.EncodeIperfServerPacket(iperfServerPacket)
	if err != nil {
		log.Fatalln("[ERROR 25] Encoding iperf server packet")
	}

	dataClient, err := t.EncodeIperfClientPacket(iperfClientPacket)
	if err != nil {
		log.Fatalln("[ERROR 23] Encoding iperf client packet")
	}

	ack.SendPacketAndWaitForAck(iperfServerPacket.PacketID, utils.SERVERID, packetsWaitingAck, &pMutex, agentSConn, nil, dataServer, "[SERVER] [MAIN READ THREAD] Iperf server packet sent", "[SERVER] [ERROR 33] Unable to send iperf server packet")
	ack.SendPacketAndWaitForAck(iperfClientPacket.PacketID, utils.SERVERID, packetsWaitingAck, &pMutex, agentCConn, nil, dataClient, "[SERVER] [MAIN READ THREAD] Iperf client packet sent", "[SERVER] [ERROR 34] Unable to send iperf client packet")
}
