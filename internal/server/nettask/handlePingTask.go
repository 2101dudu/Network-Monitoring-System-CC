package nettask

import (
	parse "nms/internal/jsonParse"
	"nms/internal/packet/ack"
	t "nms/internal/packet/task"
	utils "nms/internal/utils"
)

func handlePingTask(task parse.Task) {
	agentIP := getAgentIP(task.Devices[0].DeviceID)
	agentConn := utils.ResolveUDPAddrAndDial(agentIP, "9091")

	pingPacket := ConvertTaskIntoPingPacket(task)
	clientHash := t.CreateHashPingPacket(pingPacket)
	pingPacket.Hash = (string(clientHash))

	data := t.EncodePingPacket(pingPacket)

	ack.SendPacketAndWaitForAck(pingPacket.PacketID, utils.SERVERID, packetsWaitingAck, &pMutex, agentConn, nil, data, "[NetTask] Ping packet sent", "[ERROR 31] Unable to send ping packet")
}
