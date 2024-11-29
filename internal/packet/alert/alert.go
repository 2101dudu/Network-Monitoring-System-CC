package alert

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	utils "nms/internal/utils"
)

// Alert struct
type Alert struct { // + type of message
	PacketID   byte
	SenderID   byte
	TaskID     uint16
	CpuAlert   bool
	RamAlert   bool
	Jitter     bool
	PacketLoss bool
	Error      bool
}

// AlertBuilder struct
type AlertBuilder struct {
	Alert Alert
}

// NewAlertBuilder initializes a new AlertBuilder with default values
func NewAlertBuilder() *AlertBuilder {
	return &AlertBuilder{
		Alert: Alert{
			PacketID:   0,
			SenderID:   0,
			TaskID:     0,
			CpuAlert:   false,
			RamAlert:   false,
			Jitter:     false,
			PacketLoss: false,
			Error:      false,
		},
	}
}

// Builder methods
func (b *AlertBuilder) SetPacketID(id byte) *AlertBuilder {
	b.Alert.PacketID = id
	return b
}

func (b *AlertBuilder) SetSenderID(id byte) *AlertBuilder {
	b.Alert.SenderID = id
	return b
}

func (b *AlertBuilder) SetTaskID(id uint16) *AlertBuilder {
	b.Alert.TaskID = id
	return b
}

func (b *AlertBuilder) SetCpuAlert(alert bool) *AlertBuilder {
	b.Alert.CpuAlert = alert
	return b
}

func (b *AlertBuilder) SetRamAlert(alert bool) *AlertBuilder {
	b.Alert.RamAlert = alert
	return b
}

func (b *AlertBuilder) SetJitterAlert(alert bool) *AlertBuilder {
	b.Alert.Jitter = alert
	return b
}

func (b *AlertBuilder) SetPacketLossAlert(alert bool) *AlertBuilder {
	b.Alert.PacketLoss = alert
	return b
}

func (b *AlertBuilder) SetErrorAlert(alert bool) *AlertBuilder {
	b.Alert.Error = alert
	return b
}

func (b *AlertBuilder) Build() Alert {
	return b.Alert
}

func EncodeAlert(alert Alert) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Encode fixed fields
	buf.WriteByte(byte(utils.ALERT))
	buf.WriteByte(alert.PacketID)
	buf.WriteByte(alert.SenderID)
	// Encode TaskID
	if err := binary.Write(buf, binary.BigEndian, alert.TaskID); err != nil {
		return nil, err
	}

	// Encode boolean fields
	buf.WriteByte(utils.BoolToByte(alert.CpuAlert))
	buf.WriteByte(utils.BoolToByte(alert.RamAlert))
	buf.WriteByte(utils.BoolToByte(alert.Jitter))
	buf.WriteByte(utils.BoolToByte(alert.PacketLoss))
	buf.WriteByte(utils.BoolToByte(alert.Error))

	return buf.Bytes(), nil
}

func DecodeAlert(data []byte) (Alert, error) {
	buf := bytes.NewReader(data)
	var alert Alert

	// Decode fixed fields
	if err := binary.Read(buf, binary.BigEndian, &alert.PacketID); err != nil {
		return alert, err
	}
	if err := binary.Read(buf, binary.BigEndian, &alert.SenderID); err != nil {
		return alert, err
	}
	if err := binary.Read(buf, binary.BigEndian, &alert.TaskID); err != nil {
		return alert, err
	}

	// Decode boolean fields
	var cpuAlert, ramAlert, jitter, packetLoss, errorAlert byte
	if err := binary.Read(buf, binary.BigEndian, &cpuAlert); err != nil {
		return alert, err
	}
	if err := binary.Read(buf, binary.BigEndian, &ramAlert); err != nil {
		return alert, err
	}
	if err := binary.Read(buf, binary.BigEndian, &jitter); err != nil {
		return alert, err
	}
	if err := binary.Read(buf, binary.BigEndian, &packetLoss); err != nil {
		return alert, err
	}
	if err := binary.Read(buf, binary.BigEndian, &errorAlert); err != nil {
		return alert, err
	}

	alert.CpuAlert = cpuAlert == 1
	alert.RamAlert = ramAlert == 1
	alert.Jitter = jitter == 1
	alert.PacketLoss = packetLoss == 1
	alert.Error = errorAlert == 1

	return alert, nil
}

func EncodeAndSendAlert(conn *net.TCPConn, alert Alert) {
	alertData, err := EncodeAlert(alert)

	if err != nil {
		log.Println("[TCP][ENCODE][ERROR] Unable to encode alert")
		return
	}

	utils.WriteTCP(conn, alertData, "[TCP] Alert sent successfully", "[TCP] Failed to send alert")
}
