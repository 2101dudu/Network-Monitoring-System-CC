package jsonParse

type AlertFlowConditions struct {
	CpuUsage       byte   `json:"cpu_usage"`
	RamUsage       byte   `json:"ram_usage"`
	InterfaceStats uint16 `json:"interface_stats"`
	PacketLoss     byte   `json:"packet_loss"`
	Jitter         uint16 `json:"jitter"`
}

func validateAlertFlowConditions(afc AlertFlowConditions) bool {
	if afc.CpuUsage > 100 {
		return false
	}
	if afc.RamUsage > 100 {
		return false
	}
	if afc.PacketLoss > 100 {
		return false
	}

	return true
}
