package udp

import utils "nms/pkg/utils"

func StartUDPAgent() {
	connToServer := getUDPConnection("localhost:8081")
	registerAgent(connToServer)

	connFromAgent := utils.ResolveUDPAddrAndListen()
	defer connFromAgent.Close()
}
