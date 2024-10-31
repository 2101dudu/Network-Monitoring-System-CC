package agent_config

import (
	"bufio"
	"fmt"
	"net"
	u "nms/pkg/utils"
	"os"
)

func OpenAgent() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("[ERROR 5]: Unable to connect to server", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Established connection with the server")

	//--------------------------------------

	test_reg := u.NewRegistrationBuilder().Build()

	data, err := u.EncodeRegistration(test_reg)
	if err != nil {
		fmt.Println("[ERROR 3] Unable to enconde registration request into message", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("[ERROR 4] Unable to send message", err)
	}

	fmt.Println("Message sent")

	//-------------------------------------

	buf := make([]byte, 1024)
	n, err := bufio.NewReader(conn).Read(buf)

	if err != nil {
		fmt.Println("[ERROR 6] Unable to read message", err)
	}

	reg, err := u.DecodeRegistration(buf[:n])
	if err != nil {
		fmt.Println("[ERROR 7] Unable to decode message into registration request:", err)
	}

	if reg.NewID == 0 || !reg.SenderIsServer {
		fmt.Println("[ERROR 10] Incorrect registration parameters:", reg)
		os.Exit(1)
	}

	fmt.Println("Registration request recieved:", reg)
}
