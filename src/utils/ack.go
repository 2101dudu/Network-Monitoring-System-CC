package utils

import (
	"bytes"
	"encoding/gob"
)

type RequestType byte

const (
	AgentRegistration RequestType = iota  // iota = 0
	TaskRequest                           // iota = 1
	MetricsGathering                      // iota = 2
	Error                                 // iota = 3
)

type Ack struct {
	Acknowledged    bool
	SenderID        byte // [0, 255]
	SenderIsServer  bool
	RequestID       byte // [0, 255]
	RequestType     RequestType //(AGENT_REGISTRATION, TASK_REQUEST, METRICS_GATHERING)
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
			RequestType:    Error},
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

func (a *AckBuilder) SetRequestType(request RequestType) *AckBuilder {
	a.Ack.RequestType = request
	return a
}

func (a *AckBuilder) Build() Ack {
	return a.Ack
}

func DecodeAck(message []byte) (Ack, error) {
    var ack Ack
    buffer := bytes.NewBuffer(message)
    decoder := gob.NewDecoder(buffer)
    err := decoder.Decode(&ack)
    return ack, err
}

func EncodeAck(ack Ack) ([]byte, error) {
    var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ack)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
