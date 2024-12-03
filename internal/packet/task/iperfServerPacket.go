package task

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
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
	Bandwidth           bool
	Jitter              bool
	PacketLoss          bool
	Hash                string
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
			Bandwidth:           false,
			Jitter:              false,
			PacketLoss:          false,
			Hash:                "",
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

func (b *IperfServerPacketBuilder) SetBandwidth(bandwidth bool) *IperfServerPacketBuilder {
	b.IperfServerPacket.Bandwidth = bandwidth
	return b
}

func (b *IperfServerPacketBuilder) SetJitter(jitter bool) *IperfServerPacketBuilder {
	b.IperfServerPacket.Jitter = jitter
	return b
}

func (b *IperfServerPacket) removeHash() string {
	hash := b.Hash
	b.Hash = ""
	return hash
}

func (b *IperfServerPacketBuilder) SetPacketLoss(packetLoss bool) *IperfServerPacketBuilder {
	b.IperfServerPacket.PacketLoss = packetLoss
	return b
}

func (b *IperfServerPacketBuilder) Build() IperfServerPacket {
	return b.IperfServerPacket
}

func EncodeIperfServerPacket(msg IperfServerPacket) []byte {
	buf := new(bytes.Buffer)

	// Encode fixed fields
	buf.WriteByte(byte(utils.IPERFSERVER))
	buf.WriteByte(msg.PacketID)
	buf.WriteByte(msg.AgentID)
	binary.Write(buf, binary.BigEndian, msg.TaskID)
	binary.Write(buf, binary.BigEndian, msg.Frequency)

	// Encode DeviceMetrics
	deviceMetricsBytes := EncodeDeviceMetrics(msg.DeviceMetrics)

	buf.WriteByte(byte(len(deviceMetricsBytes))) // Add size byte
	buf.Write(deviceMetricsBytes)

	// Encode AlertFlowConditions
	alertFlowConditionsBytes := EncodeAlertFlowConditions(msg.AlertFlowConditions)

	buf.WriteByte(byte(len(alertFlowConditionsBytes))) // Add size byte
	buf.Write(alertFlowConditionsBytes)

	// Encode IperfServerCommand
	cmdBytes := []byte(msg.IperfServerCommand)
	buf.WriteByte(byte(len(cmdBytes)))
	buf.Write(cmdBytes)

	buf.WriteByte(utils.BoolToByte(msg.Bandwidth))
	buf.WriteByte(utils.BoolToByte(msg.Jitter))
	buf.WriteByte(utils.BoolToByte(msg.PacketLoss))

	// Encode Hash
	hashBytes := []byte(msg.Hash)
	buf.WriteByte(byte(len(hashBytes)))
	buf.Write(hashBytes)

	packet := buf.Bytes()

	if len(packet) > utils.BUFFERSIZE {
		log.Fatalln("[ERROR 207] Packet size too large")
	}

	return packet
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

	bandwidth, err := buf.ReadByte()
	if err != nil {
		return msg, err
	}
	jitter, err := buf.ReadByte()
	if err != nil {
		return msg, err
	}
	packetLoss, err := buf.ReadByte()
	if err != nil {
		return msg, err
	}
	msg.Bandwidth = bandwidth == 1
	msg.Jitter = jitter == 1
	msg.PacketLoss = packetLoss == 1

	// Decode Hash
	var hashLen byte
	if err := binary.Read(buf, binary.BigEndian, &hashLen); err != nil {
		return msg, err
	}
	hashBytes := make([]byte, hashLen)
	if _, err := buf.Read(hashBytes); err != nil {
		return msg, err
	}
	msg.Hash = string(hashBytes)

	return msg, nil
}

func CreateHashIperfServerPacket(iperfServer IperfServerPacket) []byte {
	byteData := EncodeIperfServerPacket(iperfServer)

	hash := sha256.Sum256(byteData)

	return hash[:utils.HASHSIZE]
}

func ValidateHashIperfServerPacket(iperfServer IperfServerPacket) bool {
	beforeHash := iperfServer.removeHash()

	byteData := EncodeIperfServerPacket(iperfServer)

	afterHash := sha256.Sum256(byteData)

	return string(afterHash[:utils.HASHSIZE]) == beforeHash
}
