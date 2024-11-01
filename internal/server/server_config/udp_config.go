package server_config

import (
	"fmt"
	"net"
	m "nms/pkg/message"
	u "nms/pkg/utils"
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
		handleUDPConnection(conn)
		agentCounter++
	}
}

func handleUDPConnection(conn *net.UDPConn) {

	fmt.Println("[UDP] Waiting for data from an agent")

	// read registration request
	n, udpAddr, regData := u.ReadUDP(conn, "[UDP] Registration request received", "[UDP] [ERROR] Unable to read registration request")

	// decode registration request
	reg, err := m.DecodeRegistration(regData[1:n])
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unbale to decode registration data:", err)
		os.Exit(1)
	}

	// validate registration request
	if reg.NewID != 0 || reg.SenderIsServer {
		fmt.Println("[UDP] [ERROR] Invalid registration request parameters")
		// ****** SEND NOACK ******
	}

	// create new registration request
	newReg := m.NewRegistrationBuilder().IsServer().SetNewID(agentCounter).Build()

	// encode new registration request
	newRegData := m.EncodeRegistration(newReg)

	// send new registration request
	u.WriteUDP(conn, udpAddr, newRegData, "[UDP] New registration request sent", "[UDP] [ERROR] Unable to send new registration request")

	// ****** SEND ACK ******

}
