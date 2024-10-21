package agent

import (
	"bufio"
	"fmt"
	"net"
	a "nms/src/utils"
	"os"
)

func Open_agent() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("[ERROR 5]: Unable to connect to server", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Established connection with the server")

	buf := make([]byte, 1024) // porque não: var buf []byte?
	n, err := bufio.NewReader(conn).Read(buf)

	if err != nil {
		fmt.Println("[ERROR 6] Unable to read message", err)
	}

	ack, err := a.Decode_ack(buf[:n]) // é mesmo necessário fornecer o tamanho do buffer?
	if err != nil {
		fmt.Println("[ERROR 7] Unable to decode message:", err)
	}

	fmt.Println("Message recieved:", ack)
}
