package agent

import (
	"bufio"
	"fmt"
	"net"
	a "nms/src/utils"
	"os"
)

func OpenAgent() {
	conn, err := net.Dial("udp", "localhost:8080")
	if err != nil {
		fmt.Println("[ERROR 5]: Unable to connect to server", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Established connection with the server")

	buf := make([]byte, 1024)
	n, err := bufio.NewReader(conn).Read(buf)

	if err != nil {
		fmt.Println("[ERROR 6] Unable to read message", err)
	}

	ack, err := a.DecodeAck(buf[:n])
	if err != nil {
		fmt.Println("[ERROR 7] Unable to decode message:", err)
	}

	fmt.Println("Message recieved:", ack)
}
