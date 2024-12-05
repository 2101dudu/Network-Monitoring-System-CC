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
	IP       [4]byte
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
			IP:       [4]byte{0, 0, 0, 0},
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

func (r *RegistrationBuilder) SetIP(ip [4]byte) *RegistrationBuilder {
	r.Registration.IP = ip
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
	if len(packet) < 7 {
		return Registration{}, errors.New("invalid packet length")
	}

	reg := Registration{
		PacketID: packet[0],
		AgentID:  packet[1],
		IP:       [4]byte{packet[2], packet[3], packet[4], packet[5]},
	}

	// Decode Hash
	hashLen := packet[6]
	if len(packet) != int(7+hashLen) {
		return Registration{}, errors.New("invalid packet length")
	}
	reg.Hash = string(packet[7 : 7+hashLen])

	return reg, nil
}

func EncodeRegistration(reg Registration) []byte {
	packet := []byte{
		byte(utils.REGISTRATION),
		reg.PacketID,
		reg.AgentID,
		reg.IP[0],
		reg.IP[1],
		reg.IP[2],
		reg.IP[3],
	}

	// Encode Hash
	hashBytes := []byte(reg.Hash)
	packet = append(packet, byte(len(hashBytes)))
	packet = append(packet, hashBytes...)

	if len(packet) > utils.BUFFERSIZE {
		log.Fatalln("[ERROR 203] Packet size too large")
	}

	return packet
}

func CreateRegistrationPacket(packetID byte, ip string) (byte, []byte) {
	byteIP, err := utils.IPStringToByte(ip)
	if err != nil {
		log.Fatalln("[ERROR 2] Unable to convert IP to byte:", err)
	}

	// generate Agent ID
	agentID, err := utils.GetAgentID()
	if err != nil {
		log.Fatalln("[ERROR 3] Unable to get agent ID:", err)
	}

	// create registration request
	registration := NewRegistrationBuilder().SetPacketID(packetID).SetAgentID(agentID).SetIP(byteIP).Build()

	hash := CreateHashRegistrationPacket(registration)

	registration.Hash = (string(hash))

	// encode registration request
	regData := EncodeRegistration(registration)

	return agentID, regData
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
