package agent_config

import (
	"fmt"
	"net"
	p "nms/pkg/packet"
	"os"
)

func ConnectTCP(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("[TCP] [ERROR] Unable to connect:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// create, encode and send registration request to server

	reg := p.NewRegistrationBuilder().Build()
	regData := p.EncodeRegistration(reg)

	_, err = conn.Write(regData)
	if err != nil {
		fmt.Println("[TCP] [ERROR] Unable to send registration request:", err)
		os.Exit(1)
	}

	fmt.Println("[TCP] Registration request sent")

	// decode new registration request from server and update registration

	newRegData := make([]byte, 1024)
	n, err := conn.Read(newRegData)
	if err != nil {
		fmt.Println("[TCP] [ERROR] Unable to read data:", err)
		os.Exit(1)
	}

	newReg, err := p.DecodeRegistration(newRegData[1:n])
	if err != nil {
		fmt.Println("[TCP] [ERROR] Unable to decode new registration data:", err)
		os.Exit(1)
	}

	//if newReg.NewID == 0 || !newReg.SenderIsServer {
	//	fmt.Println("[TCP] [ERROR] Invalid registration request parameters")
	//	// send NO_ACK
	//}

	// send ACK
	fmt.Println(newReg)
}
