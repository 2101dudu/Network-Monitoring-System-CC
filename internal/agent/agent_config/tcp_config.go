package agent_config

import (
	"fmt"
	"net"
	m "nms/pkg/message"
	"os"
)

func ConnectTCP(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Erro ao conectar via TCP:", err)
		return
	}
	defer conn.Close()

	// create, encode and send registration request to server

	reg := m.NewRegistrationBuilder().Build()
	regData := m.EncodeRegistration(reg)

	_, err = conn.Write(regData)
	if err != nil {
		fmt.Println("Unable to send registration request")
	}

	fmt.Print("Registration request sent")

	// decode new registration request from server and update registration

	newRegData := make([]byte, 1024)
	n, err := conn.Read(newRegData)
	if err != nil {
		fmt.Println("Error reading TCP data:", err)
		os.Exit(1)
	}

	newReg, err := m.DecodeRegistration(newRegData[1:n])
	if err != nil {
		fmt.Println("Error decoding new registration data:", err)
	}

	if newReg.NewID == 0 || !newReg.SenderIsServer {
		fmt.Println("Invalid registration request parameters")
		os.Exit(1)
	}

	fmt.Println(newReg)
}
