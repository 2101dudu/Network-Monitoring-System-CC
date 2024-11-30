package udp

import (
	"log"
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

	data, err := t.EncodePingPacket(pingPacket)
	if err != nil {
		log.Fatalln("[ERROR 21] Encoding ping packet")
	}

	ack.SendPacketAndWaitForAck(pingPacket.PacketID, utils.SERVERID, packetsWaitingAck, &pMutex, agentConn, nil, data, "[SERVER] [MAIN READ THREAD] Ping packet sent", "[SERVER] [ERROR 31] Unable to send ping packet")
}
