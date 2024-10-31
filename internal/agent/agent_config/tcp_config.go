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
		fmt.Println("Erro ao conectar via TCP:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Digite uma mensagem TCP: ")
		message, _ := reader.ReadString('\n')
		conn.Write([]byte(message))

		reply, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("Resposta TCP: %s", reply)
	}
}
