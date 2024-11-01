package agent_config

import (
	"fmt"
	"net"
	m "nms/pkg/message"
	u "nms/pkg/utils"
	"os"
)

func getUDPConnection(serverAddr string) *net.UDPConn {
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
	return conn
}

func ConnectUDP(serverAddr string) {
	conn := getUDPConnection(serverAddr)

	defer conn.Close()

	// create registration request
	reg := m.NewRegistrationBuilder().Build()

	// encode registration request
	regData := m.EncodeRegistration(reg)

	// send registration request
	u.WriteUDP(conn, nil, regData, "[UDP] Registration request sent", "[UDP] [ERROR] Unable to send registration request")

	// read new registration request
	n, _, newRegData := u.ReadUDP(conn, "[UDP] New registration request received", "[UDP] [ERROR] Unable to read new registration request")

	// decode new registration request (ignore the header, for now)
	newReg, err := m.DecodeRegistration(newRegData[1:n])
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to decode new registration request:", err)
		os.Exit(1)
	}

	// validate new registration request
	if newReg.NewID == 0 || !newReg.SenderIsServer {
		fmt.Println("[UDP] [ERROR] Invalid registration request parameters")
		// ****** SEND NOACK ******
	}

	// ****** SEND ACK ******
	fmt.Println(newReg)

}
