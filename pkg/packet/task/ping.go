package packet

import (
	"bytes"
	"encoding/binary"
	u "nms/pkg/utils"
)

// ------------------------- Ping ----------------------------

// -c value (packet count)
// -i value (frequency)
// ping -c 4 -i 0.5 <destination>

type PingMessage struct {
	SenderID            byte
	PacketID            byte
	TaskID              byte
	Frequency           byte
	DeviceMetrics       DeviceMetrics
	AlertFlowConditions AlertFlowConditions
	PingCommand         string
}

// ------------- Builder ----------------------
type PingMessageBuilder struct {
	pingMessage PingMessage
}

func NewPingMessageBuilder() *PingMessageBuilder {
	return &PingMessageBuilder{
		pingMessage: PingMessage{},
	}
}

func (b *PingMessageBuilder) SetSenderID(id byte) *PingMessageBuilder {
	b.pingMessage.SenderID = id
	return b
}

func (b *PingMessageBuilder) SetPacketID(id byte) *PingMessageBuilder {
	b.pingMessage.PacketID = id
	return b
}

func (b *PingMessageBuilder) SetTaskID(id byte) *PingMessageBuilder {
	b.pingMessage.TaskID = id
	return b
}

func (b *PingMessageBuilder) SetFrequency(freq byte) *PingMessageBuilder {
	b.pingMessage.Frequency = freq
	return b
}

func (b *PingMessageBuilder) SetDeviceMetrics(metrics DeviceMetrics) *PingMessageBuilder {
	b.pingMessage.DeviceMetrics = metrics
	return b
}

func (b *PingMessageBuilder) SetAlertFlowConditions(conditions AlertFlowConditions) *PingMessageBuilder {
	b.pingMessage.AlertFlowConditions = conditions
	return b
}

func (b *PingMessageBuilder) SetPingCommand(command string) *PingMessageBuilder {
	b.pingMessage.PingCommand = command
	return b
}

func (b *PingMessageBuilder) Build() PingMessage {
	return b.pingMessage
}

// -------------- Encode and Decode ---------------
func EncodePingMessage(msg PingMessage) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Encode fixed fields
	fields := []interface{}{
		byte(u.TASK),
		byte(PING),
		msg.SenderID,
		msg.PacketID,
		msg.TaskID,
		msg.Frequency,
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
	for _, interfaceString := range msg.DeviceMetrics.InterfaceStats { // Encode each interface as bytes
		interfaceBytes := []byte(interfaceString)
		err := binary.Write(buf, binary.BigEndian, byte(len(interfaceString))) // Write string length
		if err != nil {
			return nil, err
		}
		buf.Write(interfaceBytes) // Write string content
	}

	// Encode PingCommand
	cmdBytes := []byte(msg.PingCommand) // Convert string to bytes
	if err := binary.Write(buf, binary.BigEndian, byte(len(cmdBytes))); err != nil {
		return nil, err
	}
	buf.Write(cmdBytes) // Write command content

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

	// Decode PingCommand
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
