package udp

import (
	utils "nms/pkg/utils"
	"time"
)

func StartUDPAgent() {
	ip := utils.GetIPAddress()

	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")
	registerAgent(serverConn, ip)

	//agentConn := utils.ResolveUDPAddrAndListen(ip, "9091")

	//sleep for 10 seconds
	time.Sleep(10 * time.Second)
}
