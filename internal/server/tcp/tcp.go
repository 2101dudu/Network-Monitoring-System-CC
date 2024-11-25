package server

import (
	"log"
	"net"
	"os"
)

func StartTCPServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("[TCP] [ERROR] Unable to initialize the server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	log.Println("[TCP] Server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("[TCP] [ERROR] Unable to accept connection:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

// Função para tratar conexões TCP
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("[TCP] Established connection with Agent", conn.RemoteAddr())

	// decode and process registration request from agent

	regData := make([]byte, 1024)
	_, err := conn.Read(regData)
	if err != nil {
		log.Println("[TCP] [ERROR] Unable to read data:", err)
		os.Exit(1)
	}

	//reg, err := p.DecodeRegistration(regData[1:n])
	//if err != nil {
	//	log.Println("[TCP] [ERROR] Unable to decode registration data:", err)
	//	os.Exit(1)
	//}

	//if reg.NewID != 0 || reg.SenderIsServer {
	//	log.Println("[TCP] [ERROR] Invalid registration request parameters")
	//	// send NO_ACK
	//}

	// create, encode and send new registration request to agent

	//newReg := p.NewRegistrationBuilder().IsServer().SetNewID(1).Build()
	//newRegData := p.EncodeRegistration(newReg)

	//_, err = conn.Write(newRegData)
	//if err != nil {
	//	log.Println("[TCP] [ERROR] Unable to send new registration request", err)
	//	os.Exit(1)
	//}

	// send ACK
	log.Println("[TCP] New registration request sent")

}
