package agent

import (
	"log"
	"net"
	registration "nms/internal/packet/registration"
	"nms/internal/utils"
)

func ConnectTCP(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalln("[TCP] [ERROR] Unable to connect:", err)
	}
	defer conn.Close()

	// create, encode and send registration request to server
	reg := registration.NewRegistrationBuilder().Build()
	regData := registration.EncodeRegistration(reg)

	_, err = conn.Write(regData)
	if err != nil {
		log.Fatalln("[TCP] [ERROR] Unable to send registration request:", err)
	}

	log.Println("[TCP] Registration request sent")

	// decode new registration request from server and update registration
	newRegData := make([]byte, utils.BUFFERSIZE)
	n, err := conn.Read(newRegData)
	if err != nil {
		log.Fatalln("[TCP] [ERROR] Unable to read data:", err)
	}

	newReg, err := registration.DecodeRegistration(newRegData[1:n])
	if err != nil {
		log.Fatalln("[TCP] [ERROR] Unable to decode new registration data:", err)
	}

	//if newReg.NewID == 0 || !newReg.SenderIsServer {
	//	log.Println("[TCP] [ERROR] Invalid registration request parameters")
	//	// send NO_ACK
	//}

	// send ACK
	log.Println(newReg)
}
