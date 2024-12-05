package registration

import (
	"crypto/sha256"
	"errors"
	"log"
	utils "nms/internal/utils"
)

type Registration struct {
	PacketID byte
	AgentID  byte // [0, 255]
	Hash     string
}

type RegistrationBuilder struct {
	Registration Registration
}

func NewRegistrationBuilder() *RegistrationBuilder {
	return &RegistrationBuilder{
		Registration: Registration{
			PacketID: 0,
			AgentID:  0,
			Hash:     "",
		},
	}
}

func (r *RegistrationBuilder) SetPacketID(packetID byte) *RegistrationBuilder {
	r.Registration.PacketID = packetID
	return r
}

func (r *RegistrationBuilder) SetAgentID(id byte) *RegistrationBuilder {
	r.Registration.AgentID = id
	return r
}

func (r *Registration) removeHash() string {
	hash := r.Hash
	r.Hash = ""
	return hash
}

func (r *RegistrationBuilder) Build() Registration {
	return r.Registration
}

// receives the data without the header
func DecodeRegistration(packet []byte) (Registration, error) {
	if len(packet) < 2 {
		return Registration{}, errors.New("invalid packet length")
	}

	reg := Registration{
		PacketID: packet[0],
		AgentID:  packet[1],
	}

	// Decode Hash
	hashLen := packet[2]
	if len(packet) != int(3+hashLen) {
		return Registration{}, errors.New("invalid packet length")
	}
	reg.Hash = string(packet[3 : 3+hashLen])

	return reg, nil
}

func EncodeRegistration(reg Registration) []byte {
	packet := []byte{
		byte(utils.REGISTRATION),
		reg.PacketID,
		reg.AgentID,
	}

	// Encode Hash
	hashBytes := []byte(reg.Hash)
	packet = append(packet, byte(len(hashBytes)))
	packet = append(packet, hashBytes...)

	if len(packet) > utils.BUFFERSIZE {
		log.Fatalln(utils.Red+"[ERROR 203] Packet size too large", utils.Reset)
	}

	return packet
}

func CreateRegistrationPacket(packetID byte, agentID byte) []byte {
	// create registration request
	registration := NewRegistrationBuilder().SetPacketID(packetID).SetAgentID(agentID).Build()

	hash := CreateHashRegistrationPacket(registration)

	registration.Hash = (string(hash))

	// encode registration request
	regData := EncodeRegistration(registration)

	return regData
}

func CreateHashRegistrationPacket(reg Registration) []byte {
	byteData := EncodeRegistration(reg)

	hash := sha256.Sum256(byteData)

	return hash[:utils.HASHSIZE]
}

func ValidateHashRegistrationPacket(reg Registration) bool {
	beforeHash := reg.removeHash()

	byteData := EncodeRegistration(reg)

	afterHash := sha256.Sum256(byteData)

	return string(afterHash[:utils.HASHSIZE]) == beforeHash
}
