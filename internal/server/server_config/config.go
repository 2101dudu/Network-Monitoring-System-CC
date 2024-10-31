package server_config

import (
	"fmt"
	"net"
	u "nms/pkg/utils"
	"os"
)

func OpenServer(tasksList []u.Task) {

	// ---------- Uso temporário da lista de tasks ----------
	for _, task := range tasksList {
		fmt.Println("Task: ", task)
	}
	// ------------------------------------------------------

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("[ERROR 1] Uninitialized server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080...")

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

	test_ack := u.NewAckBuilder().HasAcknowledged().IsServer().SetSenderId(0).Build()

	data, err := u.EncodeAck(test_ack)
	if err != nil {
		fmt.Println("[ERROR 3] Unable to enconde message", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("[ERROR 4] Unable to send message", err)
	}

	fmt.Println("Message sent")
}
