package jsonParse

type Task struct {
	TaskID    uint16   `json:"task_id"`   // "Task-202" -> 202
	Frequency uint16   `json:"frequency"` // 0 - 65535
	Devices   []Device `json:"devices"`
}

func validateTask(task Task) bool {
	numberOfDevices := len(task.Devices)
	// check if there are more than 2 devices or no devices
	if numberOfDevices > 2 || numberOfDevices == 0 {
		return false
	}

	// iperf
	if numberOfDevices == 2 {
		// check if they are both servers or both clients
		if task.Devices[0].LinkMetrics.IperfParameters.IsServer == task.Devices[1].LinkMetrics.IperfParameters.IsServer {
			return false
		}

		// check if test durations are different
		if task.Devices[0].LinkMetrics.IperfParameters.TestDuration != task.Devices[1].LinkMetrics.IperfParameters.TestDuration {
			return false
		}

		// check if frequency is less than test duration
		if task.Frequency < task.Devices[0].LinkMetrics.IperfParameters.TestDuration {
			return false
		}

		if (task.Devices[0].LinkMetrics.IperfParameters.Bandwidth != task.Devices[1].LinkMetrics.IperfParameters.Bandwidth) || (task.Devices[0].LinkMetrics.IperfParameters.Jitter != task.Devices[1].LinkMetrics.IperfParameters.Jitter) || (task.Devices[0].LinkMetrics.IperfParameters.PacketLoss != task.Devices[1].LinkMetrics.IperfParameters.PacketLoss) {
			return false
		}
	}

	// ping
	if numberOfDevices == 1 {
		// check if ping is not enabled
		if !task.Devices[0].LinkMetrics.PingParameters.Enabled {
			return false
		}

		// check if packet count * frequency is more than metric gathering frequency
		if float32(task.Frequency) < float32(task.Devices[0].LinkMetrics.PingParameters.PacketCount)*task.Devices[0].LinkMetrics.PingParameters.Frequency {
			return false
		}
	}

	for _, device := range task.Devices {
		if !validateDevice(device, numberOfDevices) {
			return false
		}
	}

	return true
}
