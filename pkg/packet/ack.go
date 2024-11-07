package packet

import (
	"errors"
	"net"
	u "nms/pkg/utils"
)

type Ack struct {
	PacketID     byte // [0, 255]
	SenderID     byte // [0, 255]
	Acknowledged bool
}

type AckBuilder struct {
	Ack Ack
}

func NewAckBuilder() *AckBuilder {
	return &AckBuilder{
		Ack: Ack{
			PacketID:     0,
			SenderID:     0,
			Acknowledged: false},
	}
}

func (a *AckBuilder) SetPacketID(id byte) *AckBuilder {
	a.Ack.PacketID = id
	return a
}

func (a *AckBuilder) SetSenderID(id byte) *AckBuilder {
	a.Ack.SenderID = id
	return a
}

func (a *AckBuilder) HasAcknowledged() *AckBuilder {
	a.Ack.Acknowledged = true
	return a
}

func (a *AckBuilder) Build() Ack {
	return a.Ack
}

// receives the data without the header
func DecodeAck(message []byte) (Ack, error) {
	if len(message) != 3 {
		return Ack{}, errors.New("invalid message length")
	}

	ack := Ack{
		PacketID:     message[0],
		SenderID:     message[1],
		Acknowledged: message[2] == 1,
	}

	return ack, nil
}

// receives the data the header
func EncodeAck(ack Ack) []byte {
	return []byte{
		byte(u.ACK),
		ack.PacketID,
		ack.SenderID,
		u.BoolToByte(ack.Acknowledged),
	}
}

func SendAck(conn *net.UDPConn, udpAddr *net.UDPAddr, packetID byte, senderId byte, ack bool) {
	b := NewAckBuilder()
	b.SetPacketID(packetID)
	b.SetSenderID(senderId)
	if ack {
		b.HasAcknowledged()
	}

	// encode ack
	ackData := EncodeAck(b.Build())

	// send registration request
	u.WriteUDP(conn, udpAddr, ackData, "[UDP] Message sent", "[UDP] [ERROR] Unable to send message")
}
