package task

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	utils "nms/internal/utils"
)

// ------------------------- Iperf Client -----------------------------

// iperf3 -c <ServerIP>     -> Just Bandwidth (TCP adjusts automatically, no packet loss measurement)
// iperf3 -c <ServerIP> -u  -> Bandwidth, Packet Loss, and Jitter (UDP with a fixed rate set by 10M by default)

type IperfClientPacket struct {
	PacketID            byte
	AgentID             byte
	TaskID              uint16
	Frequency           uint16
	DeviceMetrics       DeviceMetrics
	AlertFlowConditions AlertFlowConditions
	IperfClientCommand  string
	Bandwidth           bool
	Jitter              bool
	PacketLoss          bool
	Hash                string
}

type IperfClientPacketBuilder struct {
	IperfClientPacket IperfClientPacket
}

func NewIperfClientPacketBuilder() *IperfClientPacketBuilder {
	return &IperfClientPacketBuilder{
		IperfClientPacket: IperfClientPacket{
			PacketID:            0,
			AgentID:             0,
			TaskID:              0,
			Frequency:           0,
			DeviceMetrics:       DeviceMetrics{},
			AlertFlowConditions: AlertFlowConditions{},
			IperfClientCommand:  "",
			Bandwidth:           false,
			Jitter:              false,
			PacketLoss:          false,
		},
	}
}

func (b *IperfClientPacketBuilder) SetPacketID(id byte) *IperfClientPacketBuilder {
	b.IperfClientPacket.PacketID = id
	return b
}

func (b *IperfClientPacketBuilder) SetAgentID(id byte) *IperfClientPacketBuilder {
	b.IperfClientPacket.AgentID = id
	return b
}

func (b *IperfClientPacketBuilder) SetTaskID(id uint16) *IperfClientPacketBuilder {
	b.IperfClientPacket.TaskID = id
	return b
}

func (b *IperfClientPacketBuilder) SetFrequency(freq uint16) *IperfClientPacketBuilder {
	b.IperfClientPacket.Frequency = freq
	return b
}

func (b *IperfClientPacketBuilder) SetDeviceMetrics(metrics DeviceMetrics) *IperfClientPacketBuilder {
	b.IperfClientPacket.DeviceMetrics = metrics
	return b
}

func (b *IperfClientPacketBuilder) SetAlertFlowConditions(conditions AlertFlowConditions) *IperfClientPacketBuilder {
	b.IperfClientPacket.AlertFlowConditions = conditions
	return b
}

func (b *IperfClientPacketBuilder) SetIperfClientCommand(cmd string) *IperfClientPacketBuilder {
	b.IperfClientPacket.IperfClientCommand = cmd
	return b
}

func (b *IperfClientPacketBuilder) SetBandwidth(bandwidth bool) *IperfClientPacketBuilder {
	b.IperfClientPacket.Bandwidth = bandwidth
	return b
}

func (b *IperfClientPacketBuilder) SetJitter(jitter bool) *IperfClientPacketBuilder {
	b.IperfClientPacket.Jitter = jitter
	return b
}

func (b *IperfClientPacketBuilder) SetPacketLoss(packetLoss bool) *IperfClientPacketBuilder {
	b.IperfClientPacket.PacketLoss = packetLoss
	return b
}

func (b *IperfClientPacket) removeHash() string {
	hash := b.Hash
	b.Hash = ""
	return hash
}

func (b *IperfClientPacketBuilder) Build() IperfClientPacket {
	return b.IperfClientPacket
}

func EncodeIperfClientPacket(msg IperfClientPacket) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Encode fixed fields
	buf.WriteByte(byte(utils.IPERFCLIENT))
	buf.WriteByte(msg.PacketID)
	buf.WriteByte(msg.AgentID)
	binary.Write(buf, binary.BigEndian, msg.TaskID)
	binary.Write(buf, binary.BigEndian, msg.Frequency)
	buf.WriteByte(utils.BoolToByte(msg.Bandwidth))
	buf.WriteByte(utils.BoolToByte(msg.Jitter))
	buf.WriteByte(utils.BoolToByte(msg.PacketLoss))

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

	// Encode IperfClientCommand
	cmdBytes := []byte(msg.IperfClientCommand)
	buf.WriteByte(byte(len(cmdBytes)))
	buf.Write(cmdBytes)

	// Encode Hash
	hashBytes := []byte(msg.Hash)
	buf.WriteByte(byte(len(hashBytes)))
	buf.Write(hashBytes)

	return buf.Bytes(), nil
}

func DecodeIperfClientPacket(data []byte) (IperfClientPacket, error) {
	buf := bytes.NewReader(data)
	var msg IperfClientPacket

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
	msg.PacketID = packetID
	msg.AgentID = agentID
	msg.Bandwidth = bandwidth == 1
	msg.Jitter = jitter == 1
	msg.PacketLoss = packetLoss == 1

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

	// Decode IperfClientCommand
	var cmdLen byte
	if err := binary.Read(buf, binary.BigEndian, &cmdLen); err != nil {
		return msg, err
	}
	cmdBytes := make([]byte, cmdLen)
	if _, err := buf.Read(cmdBytes); err != nil {
		return msg, err
	}
	msg.IperfClientCommand = string(cmdBytes)

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

func CreateHashIperfClientPacket(iperfClient IperfClientPacket) []byte {
	byteData, _ := EncodeIperfClientPacket(iperfClient)

	hash := sha256.Sum256(byteData)

	return hash[:utils.HASHSIZE]
}

func ValidateHashIperfClientPacket(iperfClient IperfClientPacket) bool {
	beforeHash := iperfClient.removeHash()

	byteData, _ := EncodeIperfClientPacket(iperfClient)

	afterHash := sha256.Sum256(byteData)

	return string(afterHash[:utils.HASHSIZE]) == beforeHash
}
