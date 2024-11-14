package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	TaskID    uint16   `json:"task_id"`   // "Task-202" -> 202
	Frequency uint16   `json:"frequency"` // 0 - 65535
	Devices   []Device `json:"devices"`
}

type Device struct {
	DeviceID      uint8         `json:"device_id"` // 0 - 255
	DeviceMetrics DeviceMetrics `json:"device_metrics"`
	LinkMetrics   LinkMetrics   `json:"link_metrics"`
}

type DeviceMetrics struct {
	CpuUsage       bool     `json:"cpu_usage"`
	RamUsage       bool     `json:"ram_usage"`
	InterfaceStats []string `json:"interface_stats"`
}

type LinkMetrics struct {
	//IsServer            bool                `json:"is_server"`
	//transportType       uint8               `json:"transportType"`
	ServerIP            [4]byte             `json:"server_ip"` // [192,168,1,2]
	TestDuration        uint16              `json:"test_duration"`
	Bandwidth           bool                `json:"bandwidth"`
	Jitter              bool                `json:"jitter"`
	PacketLoss          bool                `json:"packet_loss"`
	Latency             Latency             `json:"latency"`
	AlertFlowConditions AlertFlowConditions `json:"alertflow_conditions"`
}

type Latency struct {
	Enabled     bool   `json:"enabled"`
	Destination []byte `json:"destination"` // [192,168,1,2]
	PacketCount uint16 `json:"packet_count"`
	Frequency   uint8  `json:"frequency"`
}

type AlertFlowConditions struct {
	CpuUsage       int `json:"cpu_usage"`
	RamUsage       int `json:"ram_usage"`
	InterfaceStats int `json:"interface_stats"`
	PacketLoss     int `json:"packet_loss"`
	Jitter         int `json:"jitter"`
}

func GetDataFromJson(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("[ERROR 8] Unable to read file: ", err)
		os.Exit(1)
	}

	return data
}

func ParseDataFromJson(data []byte) []Task {
	var tasks []Task
	err := json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println("[ERROR 9] Unable to parse data", err)
		os.Exit(1)
	}

	return tasks
}
