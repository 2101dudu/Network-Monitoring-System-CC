package udp

import utils "nms/pkg/utils"

func StartUDPAgent() {
	ip := utils.GetIPAddress()

	serverConn := getUDPConnection("localhost:8081")
	registerAgent(serverConn, ip)

	agentConn := utils.ResolveUDPAddrAndListen(ip)
}
