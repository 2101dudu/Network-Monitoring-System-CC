package message

import (
	u "nms/pkg/utils"
)

type Ack struct {
	Acknowledged   bool
	SenderID       byte // [0, 255]
	SenderIsServer bool
	RequestID      byte          // [0, 255]
	RequestType    u.MessageType //(AgentRegistration, TaskRequest, MetricsGathering, Error)
}

type AckBuilder struct {
	Ack Ack
}

func NewAckBuilder() *AckBuilder {
	return &AckBuilder{
		Ack: Ack{
			Acknowledged:   false,
			SenderID:       0,
			SenderIsServer: false,
			RequestID:      0,
			RequestType:    u.ERROR},
	}
}

func (a *AckBuilder) HasAcknowledged() *AckBuilder {
	a.Ack.Acknowledged = true
	return a
}

func (a *AckBuilder) SetSenderId(id byte) *AckBuilder {
	a.Ack.SenderID = id
	return a
}

func (a *AckBuilder) IsServer() *AckBuilder {
	a.Ack.SenderIsServer = true
	return a
}

func (a *AckBuilder) SetRequestID(id byte) *AckBuilder {
	a.Ack.RequestID = id
	return a
}

func (a *AckBuilder) SetRequestType(request u.MessageType) *AckBuilder {
	a.Ack.RequestType = request
	return a
}

func (a *AckBuilder) Build() Ack {
	return a.Ack
}

// receives the data without the header
func DecodeAck(message [4]byte) Ack {
	ack := Ack{
		Acknowledged:   message[0] == 1,
		SenderID:       message[1],
		SenderIsServer: message[2] == 1,
		RequestType:    u.MessageType(message[3]),
	}

	return ack
}

// receives the data the header
func EncodeAck(ack Ack) [5]byte {
	return [5]byte{
		byte(u.ACK),
		u.BoolToByte(ack.Acknowledged),
		ack.SenderID,
		u.BoolToByte(ack.SenderIsServer),
		byte(ack.RequestType),
	}
}
