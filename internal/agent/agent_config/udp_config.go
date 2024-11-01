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

	// create registration
	reg := m.NewRegistrationBuilder().Build()

	// encode registration
	regData := m.EncodeRegistration(reg)

	// send registration
	u.WriteUDP(conn, regData, "[UDP] Registration request sent", "[UDP] [ERROR] Unable to send registration request")

	// read data
	n, newRegData := u.ReadUDP(conn, "[UDP] Data sent", "[UDP] [ERROR] Unable to read data")

	// decode data (ignore the header, for now)
	newReg, err := m.DecodeRegistration(newRegData[1:n])
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to decode new registration data:", err)
		os.Exit(1)
	}

	if newReg.NewID == 0 || !newReg.SenderIsServer {
		fmt.Println("[UDP] [ERROR] Invalid registration request parameters")
		// sendNO_ACK()
	}

	// send ACK
	fmt.Println(newReg)

}
