package server_config

import (
	"bufio"
	"fmt"
	"net"
)

func StartTCPServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor TCP:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor TCP escutando na porta", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão TCP:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

// Função para tratar conexões TCP
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler dados TCP:", err)
			break
		}
		fmt.Printf("Mensagem recebida (TCP): %s", message)
		conn.Write([]byte(message))
	}
}
