package utils

type Ack struct {
	Acknowledged   bool
	SenderID       byte // [0, 255]
	SenderIsServer bool
	RequestID      byte        // [0, 255]
	RequestType    MessageType //(AgentRegistration, TaskRequest, MetricsGathering, Error)
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
			RequestType:    ERROR},
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

func (a *AckBuilder) SetRequestType(request MessageType) *AckBuilder {
	a.Ack.RequestType = request
	return a
}

func (a *AckBuilder) Build() Ack {
	return a.Ack
}

// receives the data without the header
func DecodeAck(message [4]byte) (Ack, error) {
	ack := Ack{
		Acknowledged:   message[0] == 1,         // Decode boolean from byte
		SenderID:       message[1],              // SenderID is the 3rd byte
		SenderIsServer: message[2] == 1,         // Decode boolean from byte
		RequestType:    MessageType(message[3]), // RequestType is the 5th byte
	}

	return ack, nil
}

// receives the data the header
func EncodeAck(ack Ack) [5]byte {
	return [5]byte{
		byte(ACK),
		BoolToByte(ack.Acknowledged),
		ack.SenderID,
		BoolToByte(ack.SenderIsServer),
		byte(ack.RequestType),
	}
}
