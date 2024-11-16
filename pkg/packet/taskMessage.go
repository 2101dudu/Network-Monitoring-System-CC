package packet

import (
	"bytes"
	"encoding/binary"
	u "nms/pkg/utils"
)

// TYPE - TASK

type DeviceMetrics struct {
	CpuUsage       bool
	RamUsage       bool
	InterfaceStats []string
}

// ------------------------- Iperf Server -----------------------------
type IperfServerMessage struct {
	SenderID  byte
	PacketID  byte
	TaskID    byte
	Frequency byte
	Command   string
} // ????????

// ------------------------- Iperf Client -----------------------------

// iperf3 -c <ServerIP>     -> Just Bandwidth (TCP adjusts automatically, no packet loss measurement)
// iperf3 -c <ServerIP> -u  -> Bandwidth, Packet Loss, and Jitter (UDP with a fixed rate set by 10M by default)

type IperfMessage struct {
	SenderID      byte
	PacketID      byte
	TaskID        byte
	Frequency     byte
	DeviceMetrics DeviceMetrics
	Iperf         Iperf
}

type Iperf struct {
	TestDuration byte
	Bandwidth    bool
	Jitter       bool
	PacketLoss   bool
	IperfCommand string
}

/* type IperfMessageBuilder struct {
	IperfMessage IperfMessage
}

func NewIperfMessageBuilder() *IperfMessageBuilder {
	return &IperfMessageBuilder{
		IperfMessage: IperfMessage{
			SenderID:      0,
			PacketID:      0,
			TaskID:        0,
			Frequency:     0,
			DeviceMetrics: DeviceMetrics{},
			Iperf:         Iperf{},
		},
	}
}

func (b *IperfMessageBuilder) SetSenderID(id byte) *IperfMessageBuilder {
	b.IperfMessage.SenderID = id
	return b
}

func (b *IperfMessageBuilder) SetPacketID(id byte) *IperfMessageBuilder {
	b.IperfMessage.PacketID = id
	return b
}

func (b *IperfMessageBuilder) SetTaskID(id byte) *IperfMessageBuilder {
	b.IperfMessage.TaskID = id
	return b
}

func (b *IperfMessageBuilder) SetFrequency(freq byte) *IperfMessageBuilder {
	b.IperfMessage.Frequency = freq
	return b
}

func (b *IperfMessageBuilder) SetDeviceMetrics(metrics DeviceMetrics) *IperfMessageBuilder {
	b.IperfMessage.DeviceMetrics = metrics
	return b
}

func (b *IperfMessageBuilder) SetIperf(iperf Iperf) *IperfMessageBuilder {
	b.IperfMessage.Iperf = iperf
	return b
}

func (b *IperfMessageBuilder) Build() IperfMessage {
	return b.IperfMessage
} */

type TaskType byte

const (
	IPERF TaskType = iota
	PING
)

func EncodeIperfMessage(msg IperfMessage) ([]byte, error) {

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
	}

	for _, field := range fields { // Not necessary but if we change the types of the fields, its already done
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

func DecodeIperfMessage(data []byte) (IperfMessage, error) {
	buf := bytes.NewReader(data)
	var msg IperfMessage

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

// ------------------------- Ping ----------------------------

// -c value (packet count)
// -i value (frequency)
// ping -c 4 -i 0.5 <destination>

type PingMessage struct {
	SenderID      byte
	PacketID      byte
	TaskID        byte
	Frequency     byte
	DeviceMetrics DeviceMetrics
	PingCommand   string
}

func EncodePingMessage(msg PingMessage) ([]byte, error) {
	buf := new(bytes.Buffer)

	fields := []interface{}{
		byte(u.TASK),
		byte(PING),
		msg.SenderID,
		msg.PacketID,
		msg.TaskID,
		msg.Frequency,
		u.BoolToByte(msg.DeviceMetrics.CpuUsage),
		u.BoolToByte(msg.DeviceMetrics.RamUsage),
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

	cmdBytes := []byte(msg.PingCommand) // Convert the string to bytes
	if err := binary.Write(buf, binary.BigEndian, byte(len(cmdBytes))); err != nil {
		return nil, err
	}
	buf.Write(cmdBytes) // Write the actual string bytes to the buffer

	return buf.Bytes(), nil
}

func DecodePingMessage(data []byte) (PingMessage, error) {
	buf := bytes.NewReader(data)
	var msg PingMessage

	// Decode fixed fields
	fields := []interface{}{
		&msg.SenderID,
		&msg.PacketID,
		&msg.TaskID,
		&msg.Frequency,
		&msg.DeviceMetrics.CpuUsage,
		&msg.DeviceMetrics.RamUsage,
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

	// Decode Ping Command
	var cmdLen byte
	if err := binary.Read(buf, binary.BigEndian, &cmdLen); err != nil {
		return msg, err
	}
	cmdBytes := make([]byte, cmdLen)
	if _, err := buf.Read(cmdBytes); err != nil {
		return msg, err
	}
	msg.PingCommand = string(cmdBytes)

	return msg, nil
}
