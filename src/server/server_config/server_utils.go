package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func Open_server() {
	// Escutar numa porta específica
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Servidor a escutar na porta 8080...")

	for {
		// Aceitar conexões de clientes
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handle_connection(conn) // Lidar com a conexão numa goroutine separada
	}
}

func handle_connection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Cliente conectado:", conn.RemoteAddr())

	// Ler mensagens do cliente
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler mensagem:", err)
			return
		}
		fmt.Print("Mensagem recebida: ", message)

		// Enviar uma resposta ao cliente
		_, err = conn.Write([]byte("Mensagem recebida\n"))
		if err != nil {
			fmt.Println("Erro ao enviar resposta:", err)
			return
		}
	}
}
