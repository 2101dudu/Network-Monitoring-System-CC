package main

import agent "nms/internal/agent/udp"

func main() {

	//if protocol == "tcp" {
	//	agent.ConnectTCP(serverAddr)
	//} else if protocol == "udp" {
	//	agent.ConnectUDP(serverAddr)
	//} else {
	//	log.Println("[ERROR] Unknown procotol")
	//	os.Exit(1)
	//}

	agent.StartUDPAgent()

}
