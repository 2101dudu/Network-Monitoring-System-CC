package udp

import (
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
	serverHash := t.CreateHashIperfServerPacket(iperfServerPacket)
	iperfServerPacket.Hash = (string(serverHash))

	iperfClientPacket := ConvertTaskIntoIperfClientPacket(task, clientIndex)
	clientHash := t.CreateHashIperfClientPacket(iperfClientPacket)
	iperfClientPacket.Hash = (string(clientHash))

	dataServer := t.EncodeIperfServerPacket(iperfServerPacket)

	dataClient := t.EncodeIperfClientPacket(iperfClientPacket)

	ack.SendPacketAndWaitForAck(iperfServerPacket.PacketID, utils.SERVERID, packetsWaitingAck, &pMutex, agentSConn, nil, dataServer, "[SERVER] [MAIN READ THREAD] Iperf server packet sent", "[SERVER] [ERROR 33] Unable to send iperf server packet")
	ack.SendPacketAndWaitForAck(iperfClientPacket.PacketID, utils.SERVERID, packetsWaitingAck, &pMutex, agentCConn, nil, dataClient, "[SERVER] [MAIN READ THREAD] Iperf client packet sent", "[SERVER] [ERROR 34] Unable to send iperf client packet")
}
