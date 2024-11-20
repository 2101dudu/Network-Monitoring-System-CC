package main

import agent "nms/internal/agent/udp"

func main() {

	//if protocol == "tcp" {
	//	agent.ConnectTCP(serverAddr)
	//} else if protocol == "udp" {
	//	agent.ConnectUDP(serverAddr)
	//} else {
	//	fmt.Println("[ERROR] Unknown procotol")
	//	os.Exit(1)
	//}

	agent.ConnectUDP("localhost:8081")

}
