package packet

import (
	"errors"
	u "nms/pkg/utils"
)

type Registration struct {
	PacketID byte
	AgentID  byte // [0, 255]
}

type RegistrationBuilder struct {
	Registration Registration
}

func NewRegistrationBuilder() *RegistrationBuilder {
	return &RegistrationBuilder{
		Registration: Registration{
			PacketID: 0,
			AgentID:  0},
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

func (r *RegistrationBuilder) Build() Registration {
	return r.Registration
}

// receives the data without the header
func DecodeRegistration(message []byte) (Registration, error) {
	if len(message) != 2 {
		return Registration{}, errors.New("invalid message length")
	}

	reg := Registration{
		PacketID: message[0],
		AgentID:  message[1],
	}

	return reg, nil
}

func EncodeRegistration(reg Registration) []byte {
	return []byte{
		byte(u.REGISTRATION),
		reg.PacketID,
		reg.AgentID,
	}
}
