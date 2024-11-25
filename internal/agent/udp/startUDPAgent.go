package udp

import (
	"log"
	utils "nms/internal/utils"
	"time"
)

func StartUDPAgent() {
	// incldude "| log.Lshortfile" in the log flags to include the file name and line of code in the log
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	ip := utils.GetIPAddress()

	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")
	registerAgent(serverConn, ip)

	//agentConn := utils.ResolveUDPAddrAndListen(ip, "9091")

	//sleep for 10 seconds
	time.Sleep(10 * time.Second)
}
