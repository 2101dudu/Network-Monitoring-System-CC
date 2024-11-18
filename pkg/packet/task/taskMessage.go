package packet

// TYPE - TASK
type TaskType byte

const (
	IPERF TaskType = iota
	PING
	IPERFS
)

type DeviceMetrics struct {
	CpuUsage       bool
	RamUsage       bool
	InterfaceStats []string
}

type AlertFlowConditions struct {
	CpuUsage       byte
	RamUsage       byte
	InterfaceStats uint16
	PacketLoss     byte
	Jitter         uint16
}

// ------------------ Device Metrics Builder --------------
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

func (b *DeviceMetricsBuilder) SetCpuUsage(enabled bool) *DeviceMetricsBuilder {
	b.deviceMetrics.CpuUsage = enabled
	return b
}

func (b *DeviceMetricsBuilder) SetRamUsage(enabled bool) *DeviceMetricsBuilder {
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

// ------------------ Alert Flow Builder ----------------------
type AlertFlowConditionsBuilder struct {
	alertFlowConditions AlertFlowConditions
}

func NewAlertFlowConditionsBuilder() *AlertFlowConditionsBuilder {
	return &AlertFlowConditionsBuilder{
		alertFlowConditions: AlertFlowConditions{
			CpuUsage:       90,
			RamUsage:       90,
			InterfaceStats: 5000,
			PacketLoss:     10,
			Jitter:         100,
		},
	}
}

func (b *AlertFlowConditionsBuilder) SetCpuUsage(limit byte) *AlertFlowConditionsBuilder {
	b.alertFlowConditions.CpuUsage = limit
	return b
}

func (b *AlertFlowConditionsBuilder) SetRamUsage(limit byte) *AlertFlowConditionsBuilder {
	b.alertFlowConditions.RamUsage = limit
	return b
}

func (b *AlertFlowConditionsBuilder) SetInterfaceStats(limit uint16) *AlertFlowConditionsBuilder {
	b.alertFlowConditions.InterfaceStats = limit
	return b
}

func (b *AlertFlowConditionsBuilder) SetPacketLoss(limit byte) *AlertFlowConditionsBuilder {
	b.alertFlowConditions.PacketLoss = limit
	return b
}

func (b *AlertFlowConditionsBuilder) SetJitter(limit uint16) *AlertFlowConditionsBuilder {
	b.alertFlowConditions.Jitter = limit
	return b
}

func (b *AlertFlowConditionsBuilder) Build() AlertFlowConditions {
	return b.alertFlowConditions
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
