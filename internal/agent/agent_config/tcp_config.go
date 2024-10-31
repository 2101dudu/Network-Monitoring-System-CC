package agent_config

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func ConnectTCP(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("[ERROR] Unable to connect via TCP:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Write a TCP message: ")
		message, _ := reader.ReadString('\n')
		conn.Write([]byte(message))

		reply, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("TCP response from the server: %s", reply)
	}
}
