package tcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	alertTCP "nms/internal/packet/alert"
	"os"
	"strconv"
)

type AlertsData struct {
	TaskID    string  `json:"task_id"`
	AgentID   byte    `json:"agent_id"`
	LogTime   string  `json:"log_time"`
	AlertType string  `json:"alert_type"`
	Exceeded  float32 `json:"exceeded"`
}

func handleAlert(packetPayload []byte) {
	alert, err := alertTCP.DecodeAlert(packetPayload)
	if err != nil {
		log.Fatalln("[TCP] [ERROR 302] Unable to decode alert:", err)
	}

	// Generate an alert message dynamically based on AlertType
	var alertMessage string
	switch alert.AlertType {
	case alertTCP.CPU:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded CPU usage (%.2f%%) while executing task %d", alert.SenderID, alert.Exceeded, alert.TaskID)
	case alertTCP.RAM:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded RAM usage (%.2f%%) while executing task %d", alert.SenderID, alert.Exceeded, alert.TaskID)
	case alertTCP.JITTER:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded jitter thresholds (%.2f ms) while executing task %d", alert.SenderID, alert.Exceeded, alert.TaskID)
	case alertTCP.PACKETLOSS:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded packet loss thresholds (%.2f%%) while executing task %d", alert.SenderID, alert.Exceeded, alert.TaskID)
	case alertTCP.INTERFACESTATS:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d exceeded interface stats thresholds while executing task %d", alert.SenderID, alert.TaskID)
	case alertTCP.ERROR:
		alertMessage = fmt.Sprintf("[ALERT] Agent %d encountered an error while executing task %d", alert.SenderID, alert.TaskID)
	default:
		alertMessage = fmt.Sprintf("[TCP] [ERROR] Unknown alert type received from Agent %d for task %d", alert.SenderID, alert.TaskID)
	}

	log.Println(alertMessage)

	alertData := AlertsData{ // create json alert data
		TaskID:    "task-" + strconv.Itoa(int(alert.TaskID)),
		AgentID:   alert.SenderID,
		LogTime:   alert.Time,
		AlertType: alert.AlertType.String(),
		Exceeded:  alert.Exceeded,
	}

	// open if existes or create alerts.json if dont
	file, err := os.OpenFile("output/alerts.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("[SERVER] [ERROR 400] Unable to open alerts file:", err)
	}
	defer file.Close()

	// read json
	var alertsArray []AlertsData

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&alertsArray); err != nil && err != io.EOF {
		log.Fatalln("[SERVER] [ERROR 401] Unable to decode alerts data:", err)
	}

	// append new alert to data from json
	alertsArray = append(alertsArray, alertData)

	// write back to the file
	file.Seek(0, 0)
	file.Truncate(0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Set indentation for pretty-printing
	if err := encoder.Encode(alertsArray); err != nil {
		log.Fatalln("[SERVER] [ERROR 402] Unable to encode alerts data:", err)
	}
}
