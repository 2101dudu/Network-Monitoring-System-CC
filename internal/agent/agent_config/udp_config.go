package agent_config

import (
	"fmt"
	"net"
	p "nms/pkg/packet"
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

	id, err := u.GetAgentID()
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to get agent ID:", err)
		os.Exit(1)
	}

	// create registration request
	reg := p.NewRegistrationBuilder().SetPacketID(1).SetAgentID(id).Build()

	// encode registration request
	regData := p.EncodeRegistration(reg)

	// send registration request
	u.WriteUDP(conn, nil, regData, "[UDP] Registration request sent", "[UDP] [ERROR] Unable to send registration request")

	// for cycle - to do
	handleUDPConnection(conn)

}

func handleUDPConnection(conn *net.UDPConn) {
	fmt.Println("[UDP] Waiting for response from server")

	// read message from server
	n, _, responseData := u.ReadUDP(conn, "[UDP] Response received", "[UDP] [ERROR] Unable to read response")

	// Check if data is received
	if n == 0 {
		fmt.Println("[UDP] [ERROR] No data received")
		return
	}

	// Check message type
	msgType := u.MessageType(responseData[0])
	switch msgType {
	case u.ACK:
		fmt.Println("[UDP] Acknowledgement received from server")

	case u.ERROR:
		fmt.Println("[UDP] Error message received from server")

	case u.METRICSGATHERING:
		fmt.Println("[UDP] Metrics received from server")

	default:
		fmt.Println("[UDP] [ERROR] Unknown message type received from server")
	}
}
