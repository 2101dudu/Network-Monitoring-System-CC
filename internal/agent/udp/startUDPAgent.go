package udp

import (
	"log"
	utils "nms/internal/utils"
)

func StartUDPAgent() {
	// incldude "| log.Lshortfile" in the log flags to include the file name and line of code in the log
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	ip := utils.GetIPAddress()

	// make the agent open an UDP connection via port 9091
	agentConn := utils.ResolveUDPAddrAndListen(ip, "9091")

	// make the agent connect to the server via UDP on port 8081
	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")
	registerAgent(serverConn, ip)

	for {
		data := make([]byte, 1024)
		agentConn.Read(data)
		log.Println("[AGENT] [MAIN READ THREAD] Data received from server", ":", string(data))
	}
}
