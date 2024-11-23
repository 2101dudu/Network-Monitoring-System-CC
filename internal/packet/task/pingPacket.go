package task

import (
	"bytes"
	"encoding/binary"
	utils "nms/internal/utils"
)

// ------------------------- Ping ----------------------------

// -c value (packet count)
// -i value (frequency)
// ping -c 4 -i 0.5 <destination>

type PingMessage struct {
	AgentID             byte
	PacketID            byte
	TaskID              uint16
	Frequency           uint16
	DeviceMetrics       DeviceMetrics
	AlertFlowConditions AlertFlowConditions
	PingCommand         string
}

type PingMessageBuilder struct {
	PingMessage PingMessage
}

func NewPingMessageBuilder() *PingMessageBuilder {
	return &PingMessageBuilder{
		PingMessage: PingMessage{
			PacketID:            0,
			AgentID:             0,
			TaskID:              0,
			Frequency:           0,
			DeviceMetrics:       DeviceMetrics{},
			AlertFlowConditions: AlertFlowConditions{},
			PingCommand:         "",
		},
	}
}

func (b *PingMessageBuilder) SetPacketID(id byte) *PingMessageBuilder {
	b.PingMessage.PacketID = id
	return b
}

func (b *PingMessageBuilder) SetAgentID(id byte) *PingMessageBuilder {
	b.PingMessage.AgentID = id
	return b
}

func (b *PingMessageBuilder) SetTaskID(id uint16) *PingMessageBuilder {
	b.PingMessage.TaskID = id
	return b
}

func (b *PingMessageBuilder) SetFrequency(freq uint16) *PingMessageBuilder {
	b.PingMessage.Frequency = freq
	return b
}

func (b *PingMessageBuilder) SetDeviceMetrics(metrics DeviceMetrics) *PingMessageBuilder {
	b.PingMessage.DeviceMetrics = metrics
	return b
}

func (b *PingMessageBuilder) SetAlertFlowConditions(conditions AlertFlowConditions) *PingMessageBuilder {
	b.PingMessage.AlertFlowConditions = conditions
	return b
}

func (b *PingMessageBuilder) SetPingCommand(cmd string) *PingMessageBuilder {
	b.PingMessage.PingCommand = cmd
	return b
}

func (b *PingMessageBuilder) Build() PingMessage {
	return b.PingMessage
}

func EncodePingMessage(msg PingMessage) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Encode fixed fields
	buf.WriteByte(byte(utils.PING))
	buf.WriteByte(msg.PacketID)
	buf.WriteByte(msg.AgentID)
	binary.Write(buf, binary.BigEndian, msg.TaskID)
	binary.Write(buf, binary.BigEndian, msg.Frequency)

	// Encode DeviceMetrics
	deviceMetricsBytes, err := EncodeDeviceMetrics(msg.DeviceMetrics)
	if err != nil {
		return nil, err
	}
	buf.WriteByte(byte(len(deviceMetricsBytes))) // Add size byte
	buf.Write(deviceMetricsBytes)

	// Encode AlertFlowConditions
	alertFlowConditionsBytes, err := EncodeAlertFlowConditions(msg.AlertFlowConditions)
	if err != nil {
		return nil, err
	}
	buf.WriteByte(byte(len(alertFlowConditionsBytes))) // Add size byte
	buf.Write(alertFlowConditionsBytes)

	// Encode PingCommand
	cmdBytes := []byte(msg.PingCommand)
	buf.WriteByte(byte(len(cmdBytes)))
	buf.Write(cmdBytes)

	return buf.Bytes(), nil
}

func DecodePingMessage(data []byte) (PingMessage, error) {
	buf := bytes.NewReader(data)
	var msg PingMessage

	// Decode fixed fields
	packetID, err := buf.ReadByte()
	if err != nil {
		return msg, err
	}
	agentID, err := buf.ReadByte()
	if err != nil {
		return msg, err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.TaskID); err != nil {
		return msg, err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.Frequency); err != nil {
		return msg, err
	}
	msg.PacketID = packetID
	msg.AgentID = agentID

	// Decode DeviceMetrics
	var deviceMetricsSize byte
	if err := binary.Read(buf, binary.BigEndian, &deviceMetricsSize); err != nil {
		return msg, err
	}
	deviceMetricsBytes := make([]byte, deviceMetricsSize)
	if _, err := buf.Read(deviceMetricsBytes); err != nil {
		return msg, err
	}
	deviceMetrics, err := DecodeDeviceMetrics(deviceMetricsBytes)
	if err != nil {
		return msg, err
	}
	msg.DeviceMetrics = deviceMetrics

	// Decode AlertFlowConditions
	var alertFlowConditionsSize byte
	if err := binary.Read(buf, binary.BigEndian, &alertFlowConditionsSize); err != nil {
		return msg, err
	}
	alertFlowConditionsBytes := make([]byte, alertFlowConditionsSize)
	if _, err := buf.Read(alertFlowConditionsBytes); err != nil {
		return msg, err
	}
	alertFlowConditions, err := DecodeAlertFlowConditions(alertFlowConditionsBytes)
	if err != nil {
		return msg, err
	}
	msg.AlertFlowConditions = alertFlowConditions

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
