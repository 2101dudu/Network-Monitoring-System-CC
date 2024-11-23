package task

import (
	"bytes"
	"encoding/binary"
)

type AlertFlowConditions struct {
	CpuUsage       byte
	RamUsage       byte
	InterfaceStats uint16
	PacketLoss     byte
	Jitter         uint16
}

type AlertFlowConditionsBuilder struct {
	AlertFlowConditions AlertFlowConditions
}

func NewAlertFlowConditionsBuilder() *AlertFlowConditionsBuilder {
	return &AlertFlowConditionsBuilder{
		AlertFlowConditions: AlertFlowConditions{
			CpuUsage:       0,
			RamUsage:       0,
			InterfaceStats: 0,
			PacketLoss:     0,
			Jitter:         0,
		},
	}
}

func (b *AlertFlowConditionsBuilder) SetCpuUsage(usage byte) *AlertFlowConditionsBuilder {
	b.AlertFlowConditions.CpuUsage = usage
	return b
}

func (b *AlertFlowConditionsBuilder) SetRamUsage(usage byte) *AlertFlowConditionsBuilder {
	b.AlertFlowConditions.RamUsage = usage
	return b
}

func (b *AlertFlowConditionsBuilder) SetInterfaceStats(stats uint16) *AlertFlowConditionsBuilder {
	b.AlertFlowConditions.InterfaceStats = stats
	return b
}

func (b *AlertFlowConditionsBuilder) SetPacketLoss(loss byte) *AlertFlowConditionsBuilder {
	b.AlertFlowConditions.PacketLoss = loss
	return b
}

func (b *AlertFlowConditionsBuilder) SetJitter(jitter uint16) *AlertFlowConditionsBuilder {
	b.AlertFlowConditions.Jitter = jitter
	return b
}

func (b *AlertFlowConditionsBuilder) Build() AlertFlowConditions {
	return b.AlertFlowConditions
}

func EncodeAlertFlowConditions(conditions AlertFlowConditions) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(conditions.CpuUsage)
	buf.WriteByte(conditions.RamUsage)
	binary.Write(buf, binary.BigEndian, conditions.InterfaceStats)
	buf.WriteByte(conditions.PacketLoss)
	binary.Write(buf, binary.BigEndian, conditions.Jitter)
	return buf.Bytes(), nil
}

func DecodeAlertFlowConditions(data []byte) (AlertFlowConditions, error) {
	buf := bytes.NewReader(data)
	var conditions AlertFlowConditions

	cpuUsage, err := buf.ReadByte()
	if err != nil {
		return conditions, err
	}
	ramUsage, err := buf.ReadByte()
	if err != nil {
		return conditions, err
	}
	if err := binary.Read(buf, binary.BigEndian, &conditions.InterfaceStats); err != nil {
		return conditions, err
	}
	packetLoss, err := buf.ReadByte()
	if err != nil {
		return conditions, err
	}
	if err := binary.Read(buf, binary.BigEndian, &conditions.Jitter); err != nil {
		return conditions, err
	}
	conditions.CpuUsage = cpuUsage
	conditions.RamUsage = ramUsage
	conditions.PacketLoss = packetLoss

	return conditions, nil
}
