package server

import (
	"fmt"
	"net"
	"os"
    a "nms/src/utils"
)

func Open_server() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("[ERROR 1] Uninitialized server:", err)
		os.Exit(1)
	}
    defer listener.Close()

	fmt.Println("Server listenning on port 8080...")

    conn, err := listener.Accept()
    if err != nil {
        fmt.Println("[ERROR 2] Unable to accept connection:", err)
		os.Exit(1)
    }

    handle_connection(conn)
}

func handle_connection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Established connection with an Agent", conn.RemoteAddr())

    test_ack := a.New_ack_builder().Has_ackowledged().Is_server().Set_sender_id(0).Build()

    data, err := a.Encode_ack(test_ack)
    if err != nil {
        fmt.Println("[ERROR 3] Unable to enconde message", err)
        return
    }
    _, err = conn.Write(data) 
    fmt.Println("Message sent")
}
