package alert

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	utils "nms/internal/utils"
)

// Alert struct
type Alert struct {
	PacketID  byte
	SenderID  byte
	TaskID    uint16
	AlertType AlertType
}

// AlertType defines the type of the alert
type AlertType byte

const (
	CPU AlertType = iota
	RAM
	JITTER
	PACKETLOSS
	ERROR
)

// AlertBuilder struct
type AlertBuilder struct {
	Alert Alert
}

// NewAlertBuilder initializes a new AlertBuilder with default values
func NewAlertBuilder() *AlertBuilder {
	return &AlertBuilder{
		Alert: Alert{
			PacketID:  0,
			SenderID:  0,
			TaskID:    0,
			AlertType: ERROR, // Default to ERROR, can be changed
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

func (b *AlertBuilder) SetAlertType(alertType AlertType) *AlertBuilder {
	b.Alert.AlertType = alertType
	return b
}

func (b *AlertBuilder) Build() Alert {
	return b.Alert
}

// EncodeAlert serializes the Alert struct into a byte slice
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

	// Encode AlertType
	buf.WriteByte(byte(alert.AlertType))

	return buf.Bytes(), nil
}

// DecodeAlert deserializes the Alert struct from a byte slice
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

	// Decode AlertType
	var alertType byte
	if err := binary.Read(buf, binary.BigEndian, &alertType); err != nil {
		return alert, err
	}
	alert.AlertType = AlertType(alertType)

	return alert, nil
}

// EncodeAndSendAlert serializes and sends the alert via a TCP connection
func EncodeAndSendAlert(conn *net.TCPConn, alert Alert) {
	alertData, err := EncodeAlert(alert)
	if err != nil {
		log.Println("[TCP][ENCODE][ERROR] Unable to encode alert")
		return
	}

	utils.WriteTCP(conn, alertData, "[TCP] Alert sent successfully", "[TCP] Failed to send alert")
}
