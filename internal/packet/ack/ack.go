package ack

import (
	"crypto/sha256"
	"errors"
	"log"
	"net"
	utils "nms/internal/utils"
	"sync"
	"time"
)

type Ack struct {
	PacketID     byte // [0, 255]
	ReceiverID   byte // [0, 255]
	Acknowledged bool
	Hash         string
}

type AckBuilder struct {
	Ack Ack
}

func NewAckBuilder() *AckBuilder {
	return &AckBuilder{
		Ack: Ack{
			PacketID:     0,
			ReceiverID:   0,
			Acknowledged: false,
			Hash:         "",
		},
	}
}

func (a *AckBuilder) SetPacketID(id byte) *AckBuilder {
	a.Ack.PacketID = id
	return a
}

func (a *AckBuilder) SetReceiverID(id byte) *AckBuilder {
	a.Ack.ReceiverID = id
	return a
}

func (a *AckBuilder) HasAcknowledged() *AckBuilder {
	a.Ack.Acknowledged = true
	return a
}

func (a *Ack) removeHash() string {
	hash := a.Hash
	a.Hash = ""
	return hash
}

func (a *AckBuilder) Build() Ack {
	return a.Ack
}

// receives the data without the header
func DecodeAck(packet []byte) (Ack, error) {
	if len(packet) < 4 {
		return Ack{}, errors.New("invalid packet length")
	}

	ack := Ack{
		PacketID:     packet[0],
		ReceiverID:   packet[1],
		Acknowledged: packet[2] == 1,
	}

	// Decode Hash
	hashLen := packet[3]
	if len(packet) != int(4+hashLen) {
		return Ack{}, errors.New("invalid packet length")
	}
	ack.Hash = string(packet[4 : 4+hashLen])

	return ack, nil
}

// receives the data the header
func EncodeAck(ack Ack) []byte {
	packet := []byte{
		byte(utils.ACK),
		ack.PacketID,
		ack.ReceiverID,
		utils.BoolToByte(ack.Acknowledged),
	}

	// Encode Hash
	hashBytes := []byte(ack.Hash)
	packet = append(packet, byte(len(hashBytes)))
	packet = append(packet, hashBytes...)

	if len(packet) > utils.BUFFERSIZE {
		log.Fatalln(utils.Red, "[ERROR 201] Packet size too large", utils.Reset)
	}

	return packet
}

func EncodeAndSendAck(conn *net.UDPConn, udpAddr *net.UDPAddr, ack Ack) {
	ackData := EncodeAck(ack)
	utils.WriteUDP(conn, udpAddr, ackData, "[ERROR 14] Unable to send ack")
}

func HandleAck(ackPayload []byte, packetsWaitingAck map[byte]bool, pMutex *sync.Mutex, senderID byte) bool {
	ack, err := DecodeAck(ackPayload)
	if err != nil {
		log.Fatalln(utils.Red, "[ERROR 15] Unable to decode Ack", utils.Reset)
	}

	if !ValidateHashAckPacket(ack) {
		log.Println(utils.Red, "[ERROR 118] Invalid hash in ack packet", utils.Reset)
		return false
	}

	_, exist := utils.GetPacketStatus(ack.PacketID, packetsWaitingAck, pMutex)

	if !exist || ack.ReceiverID != senderID {
		log.Println(utils.Red, "[ERROR 16] Invalid acknowledgement", utils.Reset)
		return false
	}

	if !ack.Acknowledged {
		utils.PacketIsWaiting(ack.PacketID, packetsWaitingAck, pMutex, false)
		log.Println("[NetTask] Sender didn't acknowledge packet", ack.PacketID)
		return false
	}

	pMutex.Lock()
	delete(packetsWaitingAck, ack.PacketID)
	pMutex.Unlock()

	return true
}

func SendPacketAndWaitForAck(packetID byte, senderID byte, packetsWaitingAck map[byte]bool, pMutex *sync.Mutex, conn *net.UDPConn, udpAddr *net.UDPAddr, packetData []byte, errorMessage string) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Variable to track retransmissions
	retransmissions := 0

	// Channel to signal the main thread to stop reading
	stopReadingChan := make(chan struct{})

	go func() {
		defer wg.Done()

		// Set the status of the packet to "not" waiting for ack, because it is yet to be sent
		pMutex.Lock()
		packetsWaitingAck[packetID] = false
		pMutex.Unlock()

		packetSentInstant := time.Now()
		i := 0
		for {
			waiting, exists := utils.GetPacketStatus(packetID, packetsWaitingAck, pMutex)

			if !exists { // Registration packet has been removed from map
				return
			}

			if !waiting || time.Since(packetSentInstant) >= utils.TIMEOUTSECONDS*time.Second {
				if i >= utils.MAXRETRANSMISSIONS {
					retransmissions = i // Update retransmissions count

					// Notify the main thread to stop reading
					close(stopReadingChan)
					return
				}

				utils.WriteUDP(conn, udpAddr, packetData, errorMessage)

				utils.PacketIsWaiting(packetID, packetsWaitingAck, pMutex, true)

				packetSentInstant = time.Now()

				i++
			}
		}
	}()

	continueReadingAck := true
	for continueReadingAck {
		select {
		case <-stopReadingChan:
			// Goroutine has finished retransmissions; stop reading
			log.Println("[NetTask] Retransmissions exhausted")
			continueReadingAck = false

		default:
			// Read packet
			n, _, data := utils.ReadUDP(conn, "[ERROR 5] Unable to read ack")

			// Check if data was received
			if n == 0 {
				log.Println(utils.Red, "[ERROR 6] No data received", utils.Reset)
				return
			}

			// Get ACK contents
			packetType := utils.PacketType(data[0])
			packetPayload := data[1:n]

			if packetType != utils.ACK {
				log.Println(utils.Red, "[ERROR 17] Unexpected packet type received", utils.Reset)
				continue
			}

			continueReadingAck = !HandleAck(packetPayload, packetsWaitingAck, pMutex, senderID)
		}
	}

	// Wait for the goroutine to finish
	wg.Wait()

	// Now it is safe to close the connection
	conn.Close()

	// Check retransmission count
	if retransmissions >= utils.MAXRETRANSMISSIONS {
		log.Fatalln(utils.Red, "[ERROR 781] Unable to send packet after maximum retransmission attempts", utils.Reset)
	}
}

func CreateHashAckPacket(ack Ack) []byte {
	byteData := EncodeAck(ack)

	hash := sha256.Sum256(byteData)

	return hash[:utils.HASHSIZE]
}

func ValidateHashAckPacket(ack Ack) bool {
	hash := ack.removeHash()

	byteData := EncodeAck(ack)

	newHash := sha256.Sum256(byteData)

	return string(newHash[:utils.HASHSIZE]) == hash
}
