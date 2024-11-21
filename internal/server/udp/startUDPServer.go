package udp

import (
	"nms/pkg/utils"
)

var agentsIPs map[byte][4]byte

func StartUDPServer(port string) {
	// Initialize the map
	agentsIPs = make(map[byte][4]byte)

	serverConn := utils.ResolveUDPAddrAndListen("localhost", "8081")
	handleRegistrations(serverConn)

	//serverConn.SetDeadline(time.Now().Add(5 * time.Second))

	//serverConn.Close()

	// go Send tasks

	// go Receive metrics
}
