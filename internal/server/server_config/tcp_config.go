package server_config

import (
	"fmt"
	"net"
	m "nms/pkg/message"
	"os"
)

func StartTCPServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("[TCP] Erro ao iniciar o servidor TCP:", err)
		return
	}
	defer listener.Close()

	fmt.Println("[TCP] Servidor TCP escutando na porta", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[TCP] Erro ao aceitar conexão TCP:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

// Função para tratar conexões TCP
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("[TCP] Established connection with Agent", conn.RemoteAddr())

	// decode and process registration request from agent

	regData := make([]byte, 1024)
	n, err := conn.Read(regData)
	if err != nil {
		fmt.Println("[TCP] Error reading TCP data:", err)
		os.Exit(1)
	}

	reg, err := m.DecodeRegistration(regData[1:n])
	if err != nil {
		fmt.Println("[TCP] Error decoding registration data:", err)
	}

	if reg.NewID != 0 || reg.SenderIsServer {
		fmt.Println("[TCP] Invalid registration request parameters")
		os.Exit(1)
	}

	// create, encode and send new registration request to agent

	newReg := m.NewRegistrationBuilder().IsServer().SetNewID(1).Build()
	newRegData := m.EncodeRegistration(newReg)

	_, err = conn.Write(newRegData)
	if err != nil {
		fmt.Println("[TCP] Unable to send new registration request")
	}

	fmt.Println("[TCP] New registration request sent")

}
