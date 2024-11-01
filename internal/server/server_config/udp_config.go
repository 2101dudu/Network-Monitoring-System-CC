package server_config

import (
	"fmt"
	"net"
	m "nms/pkg/message"
	"os"
)

func StartUDPServer(port string) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Println("Erro ao resolver endereço UDP:", err)
		return
	}

	for {
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("Erro ao iniciar o servidor UDP:", err)
			os.Exit(1)
		}

		go handleUDPConnection(*conn)
	}
}

func handleUDPConnection(conn net.UDPConn) {
	defer conn.Close()

	fmt.Println("Established connection with Agent", conn.RemoteAddr())

	// decode and process registration request from agent

	regData := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(regData)
	if err != nil {
		fmt.Println("Error reading UDP data:", err)
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
