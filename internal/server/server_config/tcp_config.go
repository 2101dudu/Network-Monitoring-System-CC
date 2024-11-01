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
		fmt.Println("[TCP] [ERROR] Unable to initialize the server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("[TCP] Server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[TCP] [ERROR] Unable to accept connection:", err)
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
		fmt.Println("[TCP] [ERROR] Unable to read data:", err)
		os.Exit(1)
	}

	reg, err := m.DecodeRegistration(regData[1:n])
	if err != nil {
		fmt.Println("[TCP] [ERROR] Unable to decode registration data:", err)
		os.Exit(1)
	}

	if reg.NewID != 0 || reg.SenderIsServer {
		fmt.Println("[TCP] [ERROR] Invalid registration request parameters")
		// send NO_ACK
	}

	// create, encode and send new registration request to agent

	newReg := m.NewRegistrationBuilder().IsServer().SetNewID(1).Build()
	newRegData := m.EncodeRegistration(newReg)

	_, err = conn.Write(newRegData)
	if err != nil {
		fmt.Println("[TCP] [ERROR] Unable to send new registration request", err)
		os.Exit(1)
	}

	// send ACK
	fmt.Println("[TCP] New registration request sent")

}
