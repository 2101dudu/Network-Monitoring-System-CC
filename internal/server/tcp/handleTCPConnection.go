package tcp

import (
	"fmt"
	"log"
	"net"
	alertTCP "nms/internal/packet/alert"
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

	alert, err := alertTCP.DecodeAlert(packetPayload)
	if err != nil {
		log.Fatalln("[TCP] [ERROR 302] Unable to decode alert:", err)
	}

	// Handle the alert dynamically based on AlertType
	var alertMessage string
	switch alert.AlertType {
	case alertTCP.CPU:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded CPU usage while executing task %d", alert.SenderID, alert.TaskID)
	case alertTCP.RAM:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded RAM usage while executing task %d", alert.SenderID, alert.TaskID)
	case alertTCP.JITTER:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded jitter thresholds while executing task %d", alert.SenderID, alert.TaskID)
	case alertTCP.PACKETLOSS:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded packet loss thresholds while executing task %d", alert.SenderID, alert.TaskID)
	case alertTCP.ERROR:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d encountered an error while executing task %d", alert.SenderID, alert.TaskID)
	default:
		alertMessage = fmt.Sprintf("[TCP] [ERROR] Unknown alert type received from Agent %d for task %d", alert.SenderID, alert.TaskID)
	}

	log.Println(alertMessage)
}
