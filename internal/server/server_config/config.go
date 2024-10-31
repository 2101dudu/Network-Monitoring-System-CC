package server_config

import (
	"bufio"
	"fmt"
	"net"
	u "nms/pkg/utils"
	"os"
)

func OpenServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("[ERROR 1] Uninitialized server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listenning on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[ERROR 2] Unable to accept connection:", err)
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Established connection with an Agent", conn.RemoteAddr())

	buf := make([]byte, 1024)
	n, err := bufio.NewReader(conn).Read(buf)

	if err != nil {
		fmt.Println("[ERROR 6] Unable to read message", err)
	}

	reg, err := u.DecodeRegistration(buf[:n])
	if err != nil {
		fmt.Println("[ERROR 7] Unable to decode message into registration request:", err)
	}

	fmt.Println("Registration request recieved:", reg)

	if reg.NewID != 0 || reg.SenderIsServer {
		fmt.Println("[ERROR 10] Incorrect registration parameters:", reg)
		os.Exit(1)
	}

	regitrationToSend := u.NewRegistrationBuilder().SetNewID(1).IsServer().Build()

	data, err := u.EncodeRegistration(regitrationToSend)
	if err != nil {
		fmt.Println("[ERROR 3] Unable to enconde registration request into message", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("[ERROR 4] Unable to send message", err)
	}

	fmt.Println("Message sent")
}
