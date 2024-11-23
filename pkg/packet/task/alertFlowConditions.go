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
