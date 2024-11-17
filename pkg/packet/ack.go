package packet

import (
	"errors"
	"fmt"
	"net"
	u "nms/pkg/utils"
	"sync"
	"time"
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

func EncodeAndSendAck(conn *net.UDPConn, udpAddr *net.UDPAddr, ack Ack) {
	ackData := EncodeAck(ack)
	u.WriteUDP(conn, udpAddr, ackData, "[UDP] Message sent", "[ERROR 14] Unable to send message")
}

func HandleAck(ackPayload []byte, packetsWaitingAck map[byte]bool, pMutex *sync.Mutex, senderID byte) {
	ack, err := DecodeAck(ackPayload)
	if err != nil {
		fmt.Println("[ERROR 15] Unable to decode Ack")
		return
	}

	_, exist := GetPacketStatus(ack.PacketID, packetsWaitingAck, pMutex)

	if !exist || ack.SenderID != senderID {
		fmt.Println("[ERROR 16] Invalid acknowledgement")
		return
	}

	if !ack.Acknowledged {
		PacketIsWaiting(ack.PacketID, packetsWaitingAck, pMutex, false)
		fmt.Println("[UDP] Sender didn't acknowledge packet", ack.PacketID)
		return
	}

	pMutex.Lock()
	delete(packetsWaitingAck, ack.PacketID)
	pMutex.Unlock()
	fmt.Println("[UDP] Sender acknowledged packet", ack.PacketID)
}

func SendPacketAndWaitForAck(packetID byte, packetsWaitingAck map[byte]bool, pMutex *sync.Mutex, conn *net.UDPConn, udpAddr *net.UDPAddr, packetData []byte, successMessage string, errorMessage string) {
	packetSent := time.Now()
	for {
		waiting, exists := GetPacketStatus(packetID, packetsWaitingAck, pMutex)

		if !exists { // registration packet has been removed from map
			return
		}

		if !waiting || time.Since(packetSent) >= u.TIMEOUTSECONDS*time.Second {
			u.WriteUDP(conn, udpAddr, packetData, successMessage, errorMessage)

			PacketIsWaiting(packetID, packetsWaitingAck, pMutex, true)

			packetSent = time.Now()
		}
	}
}
