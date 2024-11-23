package registration

import (
	"errors"
	"fmt"
	utils "nms/internal/utils"
	"os"
)

type Registration struct {
	PacketID byte
	AgentID  byte // [0, 255]
	IP       [4]byte
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

func (r *RegistrationBuilder) Build() Registration {
	return r.Registration
}

// receives the data without the header
func DecodeRegistration(packet []byte) (Registration, error) {
	if len(packet) != 6 {
		return Registration{}, errors.New("invalid packet length")
	}

	reg := Registration{
		PacketID: packet[0],
		AgentID:  packet[1],
		IP:       [4]byte{packet[2], packet[3], packet[4], packet[5]},
	}

	return reg, nil
}

func EncodeRegistration(reg Registration) []byte {
	return []byte{
		byte(utils.REGISTRATION),
		reg.PacketID,
		reg.AgentID,
		reg.IP[0],
		reg.IP[1],
		reg.IP[2],
		reg.IP[3],
	}
}

func CreateRegistrationPacket(ID byte, ip string) (byte, []byte) {
	byteIP := utils.IPStringToByte(ip)

	// generate Agent ID
	agentID, err := utils.GetAgentID()
	if err != nil {
		fmt.Println("[AGENT] [ERROR 3] Unable to get agent ID:", err)
		os.Exit(1)
	}

	// create registration request
	registration := NewRegistrationBuilder().SetPacketID(ID).SetAgentID(agentID).SetIP(byteIP).Build()
	// encode registration request
	regData := EncodeRegistration(registration)

	return agentID, regData
}
