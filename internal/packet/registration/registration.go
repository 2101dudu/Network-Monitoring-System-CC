package registration

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"log"
	utils "nms/internal/utils"
)

type Registration struct {
	PacketID uint16
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

func (r *RegistrationBuilder) SetPacketID(packetID uint16) *RegistrationBuilder {
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

func DecodeRegistration(packet []byte) (Registration, error) {
	buf := bytes.NewReader(packet)
	var reg Registration

	if len(packet) < 5 { // Adjusted length check
		return reg, errors.New("invalid packet length")
	}

	if err := binary.Read(buf, binary.BigEndian, &reg.PacketID); err != nil {
		return reg, err
	}

	agentID, err := buf.ReadByte()
	if err != nil {
		return reg, err
	}
	reg.AgentID = agentID

	// Decode Hash
	var hashLen byte
	if err := binary.Read(buf, binary.BigEndian, &hashLen); err != nil {
		return reg, err
	}
	hashBytes := make([]byte, hashLen)
	if _, err := buf.Read(hashBytes); err != nil {
		return reg, err
	}
	reg.Hash = string(hashBytes)

	return reg, nil
}

func EncodeRegistration(reg Registration) []byte {
	buf := new(bytes.Buffer)

	buf.WriteByte(byte(utils.REGISTRATION))

	// Encode PacketID
	binary.Write(buf, binary.BigEndian, reg.PacketID)

	buf.WriteByte(reg.AgentID)

	// Encode Hash
	hashBytes := []byte(reg.Hash)
	buf.WriteByte(byte(len(hashBytes)))
	buf.Write(hashBytes)

	packet := buf.Bytes()

	if len(packet) > utils.BUFFERSIZE {
		log.Fatalln(utils.Red+"[ERROR 203] Packet size too large", utils.Reset)
	}

	return packet
}

func CreateRegistrationPacket(packetID uint16, agentID byte) []byte {
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
