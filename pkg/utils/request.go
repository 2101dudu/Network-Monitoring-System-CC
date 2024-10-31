package utils

import (
	"bytes"
	"encoding/gob"
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

func DecodeRegistration(message []byte) (Registration, error) {
	var reg Registration
	buffer := bytes.NewBuffer(message)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&reg)
	return reg, err
}

func EncodeRegistration(reg Registration) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(reg)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
