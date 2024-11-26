package jsonParse

type PingParameters struct {
	Enabled     bool    `json:"enabled"`
	Destination string  `json:"destination"`
	PacketCount uint16  `json:"packet_count"`
	Frequency   float32 `json:"frequency"`
}

func validatePingParameters(pingParameters PingParameters) bool {
	if len(pingParameters.Destination) == 0 && len(pingParameters.Destination) <= 15 {
		return false
	}

	return true
}
