package task

import (
	"bytes"
	"encoding/binary"
	utils "nms/internal/utils"
)

// ------------------------- Iperf Server -----------------------------
type IperfServerPacket struct {
	PacketID            byte
	AgentID             byte
	TaskID              uint16
	Frequency           uint16
	DeviceMetrics       DeviceMetrics
	AlertFlowConditions AlertFlowConditions
	IperfServerCommand  string
}

type IperfServerPacketBuilder struct {
	IperfServerPacket IperfServerPacket
}

func NewIperfServerPacketBuilder() *IperfServerPacketBuilder {
	return &IperfServerPacketBuilder{
		IperfServerPacket: IperfServerPacket{
			PacketID:            0,
			AgentID:             0,
			TaskID:              0,
			Frequency:           0,
			DeviceMetrics:       DeviceMetrics{},
			AlertFlowConditions: AlertFlowConditions{},
			IperfServerCommand:  "",
		},
	}
}

func (b *IperfServerPacketBuilder) SetPacketID(id byte) *IperfServerPacketBuilder {
	b.IperfServerPacket.PacketID = id
	return b
}

func (b *IperfServerPacketBuilder) SetAgentID(id byte) *IperfServerPacketBuilder {
	b.IperfServerPacket.AgentID = id
	return b
}

func (b *IperfServerPacketBuilder) SetTaskID(id uint16) *IperfServerPacketBuilder {
	b.IperfServerPacket.TaskID = id
	return b
}

func (b *IperfServerPacketBuilder) SetFrequency(freq uint16) *IperfServerPacketBuilder {
	b.IperfServerPacket.Frequency = freq
	return b
}

func (b *IperfServerPacketBuilder) SetDeviceMetrics(metrics DeviceMetrics) *IperfServerPacketBuilder {
	b.IperfServerPacket.DeviceMetrics = metrics
	return b
}

func (b *IperfServerPacketBuilder) SetAlertFlowConditions(conditions AlertFlowConditions) *IperfServerPacketBuilder {
	b.IperfServerPacket.AlertFlowConditions = conditions
	return b
}

func (b *IperfServerPacketBuilder) SetIperfServerCommand(cmd string) *IperfServerPacketBuilder {
	b.IperfServerPacket.IperfServerCommand = cmd
	return b
}

func (b *IperfServerPacketBuilder) Build() IperfServerPacket {
	return b.IperfServerPacket
}

func EncodeIperfServerPacket(msg IperfServerPacket) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Encode fixed fields
	buf.WriteByte(byte(utils.IPERFSERVER))
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

	// Encode IperfServerCommand
	cmdBytes := []byte(msg.IperfServerCommand)
	buf.WriteByte(byte(len(cmdBytes)))
	buf.Write(cmdBytes)

	return buf.Bytes(), nil
}

func DecodeIperfServerPacket(data []byte) (IperfServerPacket, error) {
	buf := bytes.NewReader(data)
	var msg IperfServerPacket

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

	// Decode IperfServerCommand
	var cmdLen byte
	if err := binary.Read(buf, binary.BigEndian, &cmdLen); err != nil {
		return msg, err
	}
	cmdBytes := make([]byte, cmdLen)
	if _, err := buf.Read(cmdBytes); err != nil {
		return msg, err
	}
	msg.IperfServerCommand = string(cmdBytes)

	return msg, nil
}
