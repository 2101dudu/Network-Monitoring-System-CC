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
		fmt.Println("[UDP] [ERROR] Unable to resolve address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to connect:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// create, encode and send registration request to server

	reg := m.NewRegistrationBuilder().Build()
	regData := m.EncodeRegistration(reg)

	_, err = conn.Write(regData)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to send registration request:", err)
		os.Exit(1)
	}

	fmt.Println("[UDP] Registration request sent")

	// decode new registration request from server and update registration

	newRegData := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(newRegData)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to read data:", err)
		os.Exit(1)
	}

	newReg, err := m.DecodeRegistration(newRegData[1:n])
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to decode new registration data:", err)
		os.Exit(1)
	}

	if newReg.NewID == 0 || !newReg.SenderIsServer {
		fmt.Println("[UDP] [ERROR] Invalid registration request parameters")
		// send NO_ACK
	}

	// send ACK
	fmt.Println(newReg)

}
