package agent_config

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func ConnectUDP(serverAddr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Erro ao resolver endereço UDP:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Erro ao conectar via UDP:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Digite uma mensagem UDP: ")
		message, _ := reader.ReadString('\n')
		conn.Write([]byte(message))

		buffer := make([]byte, 1024)
		n, _, _ := conn.ReadFromUDP(buffer)
		fmt.Printf("Resposta UDP: %s\n", string(buffer[:n]))
	}
}
