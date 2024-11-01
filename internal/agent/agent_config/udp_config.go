package agent_config

import (
	"fmt"
	"net"
	m "nms/pkg/message"
	"os"
)

func ConnectUDP(serverAddr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("[UDP] Erro ao resolver endereço UDP:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("[UDP] Erro ao conectar via UDP:", err)
		return
	}
	defer conn.Close()

	// create, encode and send registration request to server

	reg := m.NewRegistrationBuilder().Build()
	regData := m.EncodeRegistration(reg)

	_, err = conn.Write(regData)
	if err != nil {
		fmt.Println("[UDP] Unable to send registration request")
	}

	fmt.Println("[UDP] Registration request sent")

	// decode new registration request from server and update registration

	newRegData := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(newRegData)
	if err != nil {
		fmt.Println("[UDP] Error reading UDP data:", err)
		os.Exit(1)
	}

	newReg, err := m.DecodeRegistration(newRegData[1:n])
	if err != nil {
		fmt.Println("[UDP] Error decoding new registration data:", err)
	}

	if newReg.NewID == 0 || !newReg.SenderIsServer {
		fmt.Println("[UDP] Invalid registration request parameters")
		os.Exit(1)
	}

	fmt.Println(newReg)

}
