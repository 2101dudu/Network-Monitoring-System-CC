package main

import server "nms/internal/server"

func main() {
	//go sc.StartTCPServer("8080")
	server.StartUDPServer("8081")

	select {} // Bloqueio para manter o servidor rodando
}
