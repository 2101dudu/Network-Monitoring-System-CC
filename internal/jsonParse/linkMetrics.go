package jsonParse

type LinkMetrics struct {
	//IsServer            bool                `json:"is_server"`
	//transportType       byte               `json:"transportType"`
	IperfParameters IperfParameters `json:"iperf_parameters"`
	PingParameters  PingParameters  `json:"ping_parameters"`
}

func validateLinkMetrics(linkMetrics LinkMetrics, numberOfDevices int) bool {
	// either ping or iperf should be enabled
	if (linkMetrics.IperfParameters.Bandwidth || linkMetrics.IperfParameters.Jitter || linkMetrics.IperfParameters.PacketLoss) && linkMetrics.PingParameters.Enabled {
		return false
	}

	// if none of the iperf or ping parameters are enabled
	if !linkMetrics.IperfParameters.Bandwidth && !linkMetrics.IperfParameters.Jitter && !linkMetrics.IperfParameters.PacketLoss && !linkMetrics.PingParameters.Enabled {
		return false
	}

	// check if iperf and ping parameters are valid
	if numberOfDevices == 2 && !validateIperfParameters(linkMetrics.IperfParameters) {
		return false
	} else if numberOfDevices == 1 && !validatePingParameters(linkMetrics.PingParameters) {
		return false
	}

	return true
}
