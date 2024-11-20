package task

type AlertFlowConditions struct {
	CpuUsage       byte
	RamUsage       byte
	InterfaceStats uint16
	PacketLoss     byte
	Jitter         uint16
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
