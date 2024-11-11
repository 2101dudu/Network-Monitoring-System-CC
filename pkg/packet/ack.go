package packet

import (
	"errors"
	"fmt"
	"net"
	u "nms/pkg/utils"
	"sync"
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

func EncondeAndSendAck(conn *net.UDPConn, udpAddr *net.UDPAddr, ack Ack) {
	// encode ack
	ackData := EncodeAck(ack)

	// send registration request
	u.WriteUDP(conn, udpAddr, ackData, "[UDP] Message sent", "[UDP] [ERROR] Unable to send message")
}

func HandleAck(ackPayload []byte, packetsWaitingAck map[byte]bool, pMutex *sync.Mutex, agentID byte) {
	ack, err := DecodeAck(ackPayload)
	if err != nil {
		fmt.Println("[UDP] [ERROR] Unable to decode Ack")
		return
	} 
    pMutex.Lock()
    _, ok := packetsWaitingAck[ack.PacketID]; 
    pMutex.Unlock()
	if !ok || ack.SenderID != agentID {
		fmt.Println("[UDP] [ERROR] Invalid acknowledgement")
		return
	}

	if !ack.Acknowledged {
        pMutex.Lock()
		packetsWaitingAck[ack.PacketID] = false
        pMutex.Unlock()
		fmt.Println("[UDP] Server didn't acknowledge packet", ack.PacketID)
    } else {
        pMutex.Lock()
        delete(packetsWaitingAck, ack.PacketID)
        pMutex.Unlock()
        fmt.Println("[UDP] Server acknowledged packet", ack.PacketID)
    }
}
