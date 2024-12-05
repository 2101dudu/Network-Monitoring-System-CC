package task

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	utils "nms/internal/utils"
)

// ------------------------- Ping ----------------------------

// -c value (packet count)
// -i value (frequency)
// ping -c 4 -i 0.5 <destination>

type PingPacket struct {
	AgentID             byte
	PacketID            byte
	TaskID              uint16
	Frequency           uint16
	DeviceMetrics       DeviceMetrics
	AlertFlowConditions AlertFlowConditions
	PingCommand         string
	Hash                string
}

type PingPacketBuilder struct {
	PingPacket PingPacket
}

func NewPingPacketBuilder() *PingPacketBuilder {
	return &PingPacketBuilder{
		PingPacket: PingPacket{
			PacketID:            0,
			AgentID:             0,
			TaskID:              0,
			Frequency:           0,
			DeviceMetrics:       DeviceMetrics{},
			AlertFlowConditions: AlertFlowConditions{},
			PingCommand:         "",
			Hash:                "",
		},
	}
}

func (b *PingPacketBuilder) SetPacketID(id byte) *PingPacketBuilder {
	b.PingPacket.PacketID = id
	return b
}

func (b *PingPacketBuilder) SetAgentID(id byte) *PingPacketBuilder {
	b.PingPacket.AgentID = id
	return b
}

func (b *PingPacketBuilder) SetTaskID(id uint16) *PingPacketBuilder {
	b.PingPacket.TaskID = id
	return b
}

func (b *PingPacketBuilder) SetFrequency(freq uint16) *PingPacketBuilder {
	b.PingPacket.Frequency = freq
	return b
}

func (b *PingPacketBuilder) SetDeviceMetrics(metrics DeviceMetrics) *PingPacketBuilder {
	b.PingPacket.DeviceMetrics = metrics
	return b
}

func (b *PingPacketBuilder) SetAlertFlowConditions(conditions AlertFlowConditions) *PingPacketBuilder {
	b.PingPacket.AlertFlowConditions = conditions
	return b
}

func (b *PingPacketBuilder) SetPingCommand(cmd string) *PingPacketBuilder {
	b.PingPacket.PingCommand = cmd
	return b
}

func (b *PingPacketBuilder) Build() PingPacket {
	return b.PingPacket
}

func (b *PingPacket) removeHash() string {
	hash := b.Hash
	b.Hash = ""
	return hash
}

func EncodePingPacket(msg PingPacket) []byte {
	buf := new(bytes.Buffer)

	// Encode fixed fields
	buf.WriteByte(byte(utils.PING))
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

	// Encode PingCommand
	cmdBytes := []byte(msg.PingCommand)
	buf.WriteByte(byte(len(cmdBytes)))
	buf.Write(cmdBytes)

	// Encode Hash
	hashBytes := []byte(msg.Hash)
	buf.WriteByte(byte(len(hashBytes)))
	buf.Write(hashBytes)

	packet := buf.Bytes()

	if len(packet) > utils.BUFFERSIZE {
		log.Fatalln(utils.Red, "[ERROR 208] Packet size too large", utils.Reset)
	}

	return packet
}

func DecodePingPacket(data []byte) (PingPacket, error) {
	buf := bytes.NewReader(data)
	var msg PingPacket

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

	cmdLen, err := buf.ReadByte()
	if err != nil {
		return msg, err
	}
	cmdBytes := make([]byte, cmdLen)
	if _, err := buf.Read(cmdBytes); err != nil {
		return msg, err
	}
	msg.PingCommand = string(cmdBytes)

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

func CreateHashPingPacket(pingPacket PingPacket) []byte {
	byteData := EncodePingPacket(pingPacket)

	hash := sha256.Sum256(byteData)

	return hash[:utils.HASHSIZE]
}

func ValidateHashPingPacket(pingPacket PingPacket) bool {
	beforeHash := pingPacket.removeHash()

	byteData := EncodePingPacket(pingPacket)

	afterHash := sha256.Sum256(byteData)

	return string(afterHash[:utils.HASHSIZE]) == beforeHash
}
