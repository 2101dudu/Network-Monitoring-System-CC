package task

type DeviceMetrics struct {
	CpuUsage       bool
	RamUsage       bool
	InterfaceStats []string
}

type DeviceMetricsBuilder struct {
	deviceMetrics DeviceMetrics
}

func NewDeviceMetricsBuilder() *DeviceMetricsBuilder {
	return &DeviceMetricsBuilder{
		deviceMetrics: DeviceMetrics{
			CpuUsage:       false,
			RamUsage:       false,
			InterfaceStats: []string{},
		},
	}
}

func (b *DeviceMetricsBuilder) TrackCpuUsage(enabled bool) *DeviceMetricsBuilder {
	b.deviceMetrics.CpuUsage = enabled
	return b
}

func (b *DeviceMetricsBuilder) TrackRamUsage(enabled bool) *DeviceMetricsBuilder {
	b.deviceMetrics.RamUsage = enabled
	return b
}

func (b *DeviceMetricsBuilder) AddInterfaceStat(interfaceName string) *DeviceMetricsBuilder {
	b.deviceMetrics.InterfaceStats = append(b.deviceMetrics.InterfaceStats, interfaceName)
	return b
}

func (b *DeviceMetricsBuilder) Build() DeviceMetrics {
	return b.deviceMetrics
}

/* func HandleTask(msgData []byte, senderID byte) error {
	msgTaskType := TaskType(msgData[0]) // Iperf/Ping
	msgPayload := msgData[1:]

	switch msgTaskType {
	case IPERFS:
		fmt.Println("[AGENT] Server Iperf Command received.")
		return nil
	case IPERF:
		fmt.Println("[AGENT] Iperf Command received.")

		task, err := DecodeIperfMessage(msgPayload)
		if err != nil {
			fmt.Println("[AGENT] [ERROR] Unable to decode Iperf message data:", err)
			return err
		}

		output, err := u.ExecuteCommand(task.Iperf.IperfCommand)

		return nil
	case PING:
		fmt.Println("[AGENT] Ping Command received.")

		task, err := DecodePingMessage(msgPayload)
		if err != nil {
			fmt.Println("[AGENT] [ERROR] Unable to decode Ping message data:", err)
			return err
		}

		output, err := u.ExecuteCommand(task.PingCommand)

		return nil
	default:
		err := errors.New("unknown task type in message data")
		fmt.Println("[AGENT] [ERROR]", err)
		return err
	}

} */
