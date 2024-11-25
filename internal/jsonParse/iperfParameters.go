package jsonParse

type IperfParameters struct {
	Bandwidth    bool   `json:"bandwidth"`
	Jitter       bool   `json:"jitter"`
	PacketLoss   bool   `json:"packet_loss"`
	IsServer     bool   `json:"is_server"`
	TestDuration uint16 `json:"test_duration"`
}

func validateIperfParameters(iperfParameters IperfParameters) bool {
	// if bandwidth is enabled, jitter and packet loss should be disabled, and vice versa
	if iperfParameters.Bandwidth && (iperfParameters.Jitter || iperfParameters.PacketLoss) {
		return false
	}

	return true
}
