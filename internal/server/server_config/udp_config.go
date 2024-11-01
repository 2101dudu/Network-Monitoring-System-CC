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
		fmt.Println("[UDP] Erro ao resolver endereço UDP:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[UDP] Erro ao iniciar o servidor UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[UDP] Servidor UDP escutando na porta", port)

	for {
		handleUDPConnection(*conn)
		agentCounter++
	}
}

func handleUDPConnection(conn net.UDPConn) {
	fmt.Println("[UDP] Waiting for connection with a new agent")

	// decode and process registration request from agent

	regData := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(regData)
	if err != nil {
		fmt.Println("[UDP] Error reading UDP data:", err)
		os.Exit(1)
	}

	reg, err := m.DecodeRegistration(regData[1:n])
	if err != nil {
		fmt.Println("[UDP] Error decoding registration data:", err)
	}

	if reg.NewID != 0 || reg.SenderIsServer {
		fmt.Println("[UDP] Invalid registration request parameters")
		os.Exit(1)
	}

	// create, encode and send new registration request to agent

	newReg := m.NewRegistrationBuilder().IsServer().SetNewID(agentCounter).Build()
	newRegData := m.EncodeRegistration(newReg)

	_, err = conn.WriteToUDP(newRegData, addr)
	if err != nil {
		fmt.Println("[UDP] Unable to send new registration request")
		os.Exit(1)
	}

	fmt.Println("[UDP] New registration request sent")

}
