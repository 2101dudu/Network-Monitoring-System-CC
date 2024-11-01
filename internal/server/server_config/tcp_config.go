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
		fmt.Println("Erro ao iniciar o servidor TCP:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor TCP escutando na porta", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão TCP:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

// Função para tratar conexões TCP
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Established connection with Agent", conn.RemoteAddr())

	// decode and process registration request from agent

	regData := make([]byte, 1024)
	n, err := conn.Read(regData)
	if err != nil {
		fmt.Println("Error reading TCP data:", err)
		os.Exit(1)
	}

	reg, err := m.DecodeRegistration(regData[:n])
	if err != nil {
		fmt.Println("Error decoding registration data:", err)
	}

	if reg.NewID != 0 || reg.SenderIsServer {
		fmt.Println("Invalid registration request parameters")
		os.Exit(1)
	}

	// create, encode and send new registration request to agent

	newReg := m.NewRegistrationBuilder().IsServer().SetNewID(1).Build()
	newRegData := m.EncodeRegistration(newReg)

	_, err = conn.Write(newRegData)
	if err != nil {
		fmt.Println("Unable to send new registration request")
	}

	fmt.Print("New registration request sent")

}
