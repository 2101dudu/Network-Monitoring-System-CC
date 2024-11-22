package jsonParse

type DeviceMetrics struct {
	CpuUsage       bool     `json:"cpu_usage"`
	RamUsage       bool     `json:"ram_usage"`
	InterfaceStats []string `json:"interface_stats"`
}

func validateDeviceMetrics(dm DeviceMetrics) bool {
	if len(dm.InterfaceStats) == 0 {
		return false
	}

	return true
}
