package server_config

import (
	"fmt"
	"net"
	m "nms/pkg/message"
	"os"
)

var agentCounter byte = 1

func StartUDPServer(port string) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to resolve address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to initialize the server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[UDP] Server listening on port", port)

	for {
		handleUDPConnection(*conn)
		agentCounter++
	}
}

func handleUDPConnection(conn net.UDPConn) {

	fmt.Println("[UDP] Waiting for data from an agent")

	regData := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(regData)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to read data:", err)
		os.Exit(1)
	}

	reg, err := m.DecodeRegistration(regData[1:n])
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unbale to decode registration data:", err)
		os.Exit(1)
	}

	if reg.NewID != 0 || reg.SenderIsServer {
		fmt.Println("[UDP] [ERROR] Invalid registration request parameters")
		// send NO_ACK
	}

	// create, encode and send new registration request to agent

	newReg := m.NewRegistrationBuilder().IsServer().SetNewID(agentCounter).Build()
	newRegData := m.EncodeRegistration(newReg)

	_, err = conn.WriteToUDP(newRegData, addr)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to send new registration request", err)
		os.Exit(1)
	}

	// send ACK
	fmt.Println("[UDP] New registration request sent")

}
