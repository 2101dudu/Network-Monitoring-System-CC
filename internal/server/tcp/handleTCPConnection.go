package tcp

import (
	"log"
	"net"
	utils "nms/internal/utils"
)

// func to handle tcps connections with agents
func handleTCPConnection(conn *net.TCPConn) {
	defer conn.Close()

	//log.Println("[TCP] Established connection with Agent ", conn.RemoteAddr())

	n, alertData := utils.ReadTCP(conn, "[TCP] Sucess reading alert data", "[TCP] [ERROR 299] Unable to read alert data")

	// Check if there is data
	if n == 0 {
		log.Println("[TCP] [ERROR 300] No data received")
		return
	}

	// type cast the data to the appropriate packet type
	packetType := utils.PacketType(alertData[0])
	packetPayload := alertData[1:n]

	if packetType != utils.ALERT {
		log.Fatalln("[TCP] [ERROR 301] Unexpected packet type received from agent")
	}

	handleAlert(packetPayload)
}
