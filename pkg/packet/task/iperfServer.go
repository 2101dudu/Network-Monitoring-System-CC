package task

import (
	"bytes"
	"encoding/binary"
	u "nms/pkg/utils"
)

// ------------------------- Iperf Server -----------------------------
type IperfServerMessage struct {
	SenderID byte
	PacketID byte
	TaskID   byte
	Command  string
}

// ---------------------- Builder -----------------------------
type IperfServerMessageBuilder struct {
	iperfServerMessage IperfServerMessage
}

func NewIperfServerMessageBuilder() *IperfServerMessageBuilder {
	return &IperfServerMessageBuilder{
		iperfServerMessage: IperfServerMessage{},
	}
}

func (b *IperfServerMessageBuilder) SetSenderID(id byte) *IperfServerMessageBuilder {
	b.iperfServerMessage.SenderID = id
	return b
}

func (b *IperfServerMessageBuilder) SetPacketID(id byte) *IperfServerMessageBuilder {
	b.iperfServerMessage.PacketID = id
	return b
}

func (b *IperfServerMessageBuilder) SetTaskID(id byte) *IperfServerMessageBuilder {
	b.iperfServerMessage.TaskID = id
	return b
}
func (b *IperfServerMessageBuilder) SetIperfServerMessageCommand(command string) *IperfServerMessageBuilder {
	b.iperfServerMessage.Command = command
	return b
}

func (b *IperfServerMessageBuilder) Build() IperfServerMessage {
	return b.iperfServerMessage
}

// ------------------- Encode and Decode --------------------------

func EncodeIperfServerMessage(msg IperfServerMessage) ([]byte, error) {
	buf := new(bytes.Buffer)

	fields := []interface{}{
		byte(u.TASK),
		byte(PING),
		msg.SenderID,
		msg.PacketID,
		msg.TaskID,
	}

	for _, field := range fields {
		err := binary.Write(buf, binary.BigEndian, field)
		if err != nil {
			return nil, err
		}
	}

	cmdBytes := []byte(msg.Command)
	if err := binary.Write(buf, binary.BigEndian, byte(len(cmdBytes))); err != nil {
		return nil, err
	}
	buf.Write(cmdBytes)

	return buf.Bytes(), nil
}

func DecodeIperfServerMessage(data []byte) (IperfServerMessage, error) {
	buf := bytes.NewReader(data)
	var msg IperfServerMessage

	fields := []interface{}{
		&msg.SenderID,
		&msg.PacketID,
		&msg.TaskID,
	}

	for _, field := range fields {
		if err := binary.Read(buf, binary.BigEndian, field); err != nil {
			return msg, err
		}
	}

	var cmdLen byte
	if err := binary.Read(buf, binary.BigEndian, &cmdLen); err != nil {
		return msg, err
	}
	cmdBytes := make([]byte, cmdLen)
	if _, err := buf.Read(cmdBytes); err != nil {
		return msg, err
	}
	msg.Command = string(cmdBytes)

	return msg, nil
}
