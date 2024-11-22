package jsonParse

type AlertFlowConditions struct {
	CpuUsage       int `json:"cpu_usage"`
	RamUsage       int `json:"ram_usage"`
	InterfaceStats int `json:"interface_stats"`
	PacketLoss     int `json:"packet_loss"`
	Jitter         int `json:"jitter"`
}

func validateAlertFlowConditions(afc AlertFlowConditions) bool {
	if afc.CpuUsage < 0 || afc.CpuUsage > 100 {
		return false
	}
	if afc.RamUsage < 0 || afc.RamUsage > 100 {
		return false
	}
	if afc.PacketLoss < 0 || afc.PacketLoss > 100 {
		return false
	}

	return true
}
