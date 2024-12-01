package tcp

import (
	"log"
	"net"
	alertTCP "nms/internal/packet/alert"
	utils "nms/internal/utils"
)

// func to handle tcps connections with agents
func handleTCPConnection(conn *net.TCPConn) {
	defer conn.Close()

	//log.Println("[TCP] Established connection with Agent ", conn.RemoteAddr())

	n, alertData := utils.ReadTCP(conn, "[TCP] Sucess reading alert data", "[TCP] [ERROR] Unable to read alert data")

	// Check if there is data
	if n == 0 {
		log.Println("[TCP] [ERROR] No data received")
		return
	}

	// type cast the data to the appropriate packet type
	packetType := utils.PacketType(alertData[0])
	packetPayload := alertData[1:n]

	if packetType != utils.ALERT {
		log.Fatalln("[TCP] [ERROR] Unexpected packet type received from agent")
	}

	alert, err := alertTCP.DecodeAlert(packetPayload)
	if err != nil {
		log.Fatalln("[TCP] [ERROR] Unable to decode alert:", err)
	}

	if alert.CpuAlert {
		log.Println("[ALERT] Agent", alert.SenderID, "exceeded CPU usage while executing task ", alert.TaskID)
	}
	if alert.RamAlert {
		log.Println("[ALERT] Agent", alert.SenderID, "exceeded RAM usage while executing task ", alert.TaskID)
	}

	//log.Println("[TCP] Closing connection with ", conn.RemoteAddr())
}
