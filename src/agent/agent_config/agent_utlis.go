package agent

import (
	"bufio"
	"fmt"
	"net"
	"os"
    a "nms/src/utils"
)

func Open_agent() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("[ERROR 4]: Unable to connect to server", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Established connection with the server")

	for {
        buf := make([]byte, 1024)   
        n, err := bufio.NewReader(conn).Read(buf)

        if err != nil {
			fmt.Println("[ERROR 5] Unable to read message", err)
        }

        ack, err := a.Decode_ack(buf[:n])
        if err != nil {
			fmt.Println("[ERRO 6] Unable to decode message:", err)
        }

        fmt.Print("Message recieved:", ack)
	}
}
