package server_config

import (
	"fmt"
	"net"
	"strings"
)

func StartUDPServer(port string) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Println("Erro ao resolver endereço UDP:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Servidor UDP escutando na porta", port)

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao ler dados UDP:", err)
			continue
		}
		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("Mensagem recebida (UDP): %s\n", message)
		conn.WriteToUDP([]byte(message), clientAddr)
	}
}
