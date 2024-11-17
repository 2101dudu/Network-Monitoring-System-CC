package packet

import (
	"bytes"
	"encoding/binary"
	u "nms/pkg/utils"
)

// ------------------------- Iperf Client -----------------------------

// iperf3 -c <ServerIP>     -> Just Bandwidth (TCP adjusts automatically, no packet loss measurement)
// iperf3 -c <ServerIP> -u  -> Bandwidth, Packet Loss, and Jitter (UDP with a fixed rate set by 10M by default)

type IperfClientMessage struct {
	SenderID            byte
	PacketID            byte
	TaskID              byte
	Frequency           byte
	DeviceMetrics       DeviceMetrics
	AlertFlowConditions AlertFlowConditions
	Iperf               Iperf
}

type Iperf struct {
	TestDuration byte
	Bandwidth    bool
	Jitter       bool
	PacketLoss   bool
	IperfCommand string
}

/* type IperfClientMessageBuilder struct {
	IperfClientMessage IperfClientMessage
}

func NewIperfClientMessageBuilder() *IperfClientMessageBuilder {
	return &IperfClientMessageBuilder{
		IperfClientMessage: IperfClientMessage{
			SenderID:      0,
			PacketID:      0,
			TaskID:        0,
			Frequency:     0,
			DeviceMetrics: DeviceMetrics{},
			Iperf:         Iperf{},
		},
	}
}

func (b *IperfClientMessageBuilder) SetSenderID(id byte) *IperfClientMessageBuilder {
	b.IperfClientMessage.SenderID = id
	return b
}

func (b *IperfClientMessageBuilder) SetPacketID(id byte) *IperfClientMessageBuilder {
	b.IperfClientMessage.PacketID = id
	return b
}

func (b *IperfClientMessageBuilder) SetTaskID(id byte) *IperfClientMessageBuilder {
	b.IperfClientMessage.TaskID = id
	return b
}

func (b *IperfClientMessageBuilder) SetFrequency(freq byte) *IperfClientMessageBuilder {
	b.IperfClientMessage.Frequency = freq
	return b
}

func (b *IperfClientMessageBuilder) SetDeviceMetrics(metrics DeviceMetrics) *IperfClientMessageBuilder {
	b.IperfClientMessage.DeviceMetrics = metrics
	return b
}

func (b *IperfClientMessageBuilder) SetIperf(iperf Iperf) *IperfClientMessageBuilder {
	b.IperfClientMessage.Iperf = iperf
	return b
}

func (b *IperfClientMessageBuilder) Build() IperfClientMessage {
	return b.IperfClientMessage
} */

func EncodeIperfClientMessage(msg IperfClientMessage) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Encode fixed fields
	fields := []interface{}{
		byte(u.TASK),
		byte(IPERF),
		msg.SenderID,
		msg.PacketID,
		msg.TaskID,
		msg.Frequency,
		u.BoolToByte(msg.Iperf.Bandwidth),
		u.BoolToByte(msg.Iperf.Jitter),
		u.BoolToByte(msg.Iperf.PacketLoss),
		msg.Iperf.TestDuration,
		u.BoolToByte(msg.DeviceMetrics.CpuUsage),
		u.BoolToByte(msg.DeviceMetrics.RamUsage),
		msg.AlertFlowConditions.CpuUsage,
		msg.AlertFlowConditions.RamUsage,
		msg.AlertFlowConditions.InterfaceStats,
		msg.AlertFlowConditions.PacketLoss,
		msg.AlertFlowConditions.Jitter,
	}

	for _, field := range fields {
		err := binary.Write(buf, binary.BigEndian, field)
		if err != nil {
			return nil, err
		}
	}

	// Encode InterfaceStats
	err := binary.Write(buf, binary.BigEndian, byte(len(msg.DeviceMetrics.InterfaceStats))) // Write length of the array of strings
	if err != nil {
		return nil, err
	}
	for _, interfaceString := range msg.DeviceMetrics.InterfaceStats { //For each string transforms to bytes and then writes each one on the buffer with their length
		interfaceBytes := []byte(interfaceString)
		err := binary.Write(buf, binary.BigEndian, byte(len(interfaceString))) // Write length of the string
		if err != nil {
			return nil, err
		}
		buf.Write(interfaceBytes) // Write the actual string bytes to the buffer
	}

	// Encode IperfCommand
	cmdBytes := []byte(msg.Iperf.IperfCommand) // Convert the string to bytes
	if err := binary.Write(buf, binary.BigEndian, byte(len(cmdBytes))); err != nil {
		return nil, err
	}
	buf.Write(cmdBytes) // Write the actual string bytes to the buffer

	return buf.Bytes(), nil
}

func DecodeIperfClientMessage(data []byte) (IperfClientMessage, error) {
	buf := bytes.NewReader(data)
	var msg IperfClientMessage

	// Decode fixed fields
	fields := []interface{}{
		&msg.SenderID,
		&msg.PacketID,
		&msg.TaskID,
		&msg.Frequency,
		&msg.Iperf.Bandwidth,
		&msg.Iperf.Jitter,
		&msg.Iperf.PacketLoss,
		&msg.Iperf.TestDuration,
		&msg.DeviceMetrics.CpuUsage,
		&msg.DeviceMetrics.RamUsage,
		&msg.AlertFlowConditions.CpuUsage,
		&msg.AlertFlowConditions.RamUsage,
		&msg.AlertFlowConditions.InterfaceStats,
		&msg.AlertFlowConditions.PacketLoss,
		&msg.AlertFlowConditions.Jitter,
	}

	for _, field := range fields {
		if err := binary.Read(buf, binary.BigEndian, field); err != nil {
			return msg, err
		}
	}

	// Decode InterfaceStats
	var interfaceCount byte
	if err := binary.Read(buf, binary.BigEndian, &interfaceCount); err != nil {
		return msg, err
	}
	msg.DeviceMetrics.InterfaceStats = make([]string, interfaceCount)
	for i := range msg.DeviceMetrics.InterfaceStats {
		var interfaceLen byte
		if err := binary.Read(buf, binary.BigEndian, &interfaceLen); err != nil {
			return msg, err
		}
		interfaceBytes := make([]byte, interfaceLen)
		if _, err := buf.Read(interfaceBytes); err != nil {
			return msg, err
		}
		msg.DeviceMetrics.InterfaceStats[i] = string(interfaceBytes)
	}

	// Decode IperfCommand
	var cmdLen byte
	if err := binary.Read(buf, binary.BigEndian, &cmdLen); err != nil {
		return msg, err
	}
	cmdBytes := make([]byte, cmdLen)
	if _, err := buf.Read(cmdBytes); err != nil {
		return msg, err
	}
	msg.Iperf.IperfCommand = string(cmdBytes)

	return msg, nil
}
