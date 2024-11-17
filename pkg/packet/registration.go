package packet

import (
	"errors"
	"fmt"
	u "nms/pkg/utils"
	"os"
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

func CreateRegistrationPacket(ID byte) (byte, []byte) {
	// generate Agent ID
	agentID, err := u.GetAgentID()
	if err != nil {
		fmt.Println("[AGENT] [ERROR 3] Unable to get agent ID:", err)
		os.Exit(1)
	}

	// create registration request
	registration := NewRegistrationBuilder().SetPacketID(ID).SetAgentID(agentID).Build()
	// encode registration request
	regData := EncodeRegistration(registration)

	return agentID, regData
}
