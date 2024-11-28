package ack

import (
	"errors"
	"log"
	"net"
	utils "nms/internal/utils"
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
func DecodeAck(packet []byte) (Ack, error) {
	if len(packet) != 3 {
		return Ack{}, errors.New("invalid packet length")
	}

	ack := Ack{
		PacketID:     packet[0],
		SenderID:     packet[1],
		Acknowledged: packet[2] == 1,
	}

	return ack, nil
}

// receives the data the header
func EncodeAck(ack Ack) []byte {
	return []byte{
		byte(utils.ACK),
		ack.PacketID,
		ack.SenderID,
		utils.BoolToByte(ack.Acknowledged),
	}
}

func EncodeAndSendAck(conn *net.UDPConn, udpAddr *net.UDPAddr, ack Ack) {
	ackData := EncodeAck(ack)
	utils.WriteUDP(conn, udpAddr, ackData, "[UDP] Ack sent", "[ERROR 14] Unable to send ack")
}

func HandleAck(ackPayload []byte, packetsWaitingAck map[byte]bool, pMutex *sync.Mutex, senderID byte, conn *net.UDPConn) bool {
	ack, err := DecodeAck(ackPayload)
	if err != nil {
		log.Println("[ERROR 15] Unable to decode Ack")
		return false
	}

	_, exist := utils.GetPacketStatus(ack.PacketID, packetsWaitingAck, pMutex)

	if !exist || ack.SenderID != senderID {
		log.Println("[ERROR 16] Invalid acknowledgement")
		return false
	}

	if !ack.Acknowledged {
		utils.PacketIsWaiting(ack.PacketID, packetsWaitingAck, pMutex, false)
		log.Println("[UDP] Sender didn't acknowledge packet", ack.PacketID)
		return false
	}

	pMutex.Lock()
	delete(packetsWaitingAck, ack.PacketID)
	pMutex.Unlock()
	log.Println("[UDP] Sender acknowledged packet", ack.PacketID)

	// if the packet was acknowledged, the connection can be closed, as it's no longer needed
	conn.Close()
	return true
}

func SendPacketAndWaitForAck(packetID byte, senderID byte, packetsWaitingAck map[byte]bool, pMutex *sync.Mutex, conn *net.UDPConn, udpAddr *net.UDPAddr, packetData []byte, successMessage string, errorMessage string) {
	// set the status of the packet to "not" waiting for ack, because it is yet to be sent
	utils.PacketIsWaiting(packetID, packetsWaitingAck, pMutex, false)

	packetSentInstant := time.Now()
	go func() {
		for {
			waiting, exists := utils.GetPacketStatus(packetID, packetsWaitingAck, pMutex)

			if !exists { // registration packet has been removed from map
				return
			}

			if !waiting || time.Since(packetSentInstant) >= utils.TIMEOUTSECONDS*time.Second {
				utils.WriteUDP(conn, udpAddr, packetData, successMessage, errorMessage)

				utils.PacketIsWaiting(packetID, packetsWaitingAck, pMutex, true)

				packetSentInstant = time.Now()
			}

			// add a small delay to prevent the loop from running too fast
			//time.Sleep(1 * time.Millisecond)
		}
	}()

	ackWasSent := false
	for !ackWasSent {
		log.Println("[UDP] Waiting for ack")

		// read packet
		n, _, data := utils.ReadUDP(conn, "[UDP] Ack received", "[UDP] [ERROR 5] Unable to read ack")

		// Check if data was received
		if n == 0 {
			log.Println("[UDP] [ERROR 6] No data received")
			return
		}

		// get ACK contents
		packetType := utils.PacketType(data[0])
		packetPayload := data[1:n]

		if packetType != utils.ACK {
			log.Println("[UDP] [ERROR 17] Unexpected packet type received")
			return
		}
		ackWasSent = HandleAck(packetPayload, packetsWaitingAck, pMutex, senderID, conn)
	}
}
