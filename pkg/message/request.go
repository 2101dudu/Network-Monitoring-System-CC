package message

import (
	u "nms/pkg/utils"
	// u "nms/pkg/utils"
)

type Registration struct {
	SenderIsServer bool
	NewID          byte // [0, 255]
}

type RegistrationBuilder struct {
	Registration Registration
}

func NewRegistrationBuilder() *RegistrationBuilder {
	return &RegistrationBuilder{
		Registration: Registration{
			SenderIsServer: false,
			NewID:          0},
	}
}

func (r *RegistrationBuilder) IsServer() *RegistrationBuilder {
	r.Registration.SenderIsServer = true
	return r
}

func (r *RegistrationBuilder) SetNewID(id byte) *RegistrationBuilder {
	r.Registration.NewID = id
	return r
}

func (r *RegistrationBuilder) Build() Registration {
	return r.Registration
}

// receives the data without the header
func DecodeRegistration(message [2]byte) Registration {
	reg := Registration{
		SenderIsServer: message[0] == 1,
		NewID:          message[1],
	}

	return reg
}

func EncodeRegistration(reg Registration) [3]byte {
	return [3]byte{
		byte(u.REGSITRATION),
		u.BoolToByte(reg.SenderIsServer),
		reg.NewID,
	}
}
