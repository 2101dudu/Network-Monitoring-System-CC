package jsonParse

type PingParameters struct {
	Enabled     bool   `json:"enabled"`
	Destination []byte `json:"destination"` // [192,168,1,2]
	PacketCount uint16 `json:"packet_count"`
	Frequency   uint8  `json:"frequency"`
}

func validatePingParameters(pingParameters PingParameters) bool {
	if len(pingParameters.Destination) != 4 {
		return false
	}

	return true
}
