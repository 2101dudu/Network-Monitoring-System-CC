package main

import sc "nms/internal/server/server_config"

func main() {
	go sc.StartTCPServer("8080")
	go sc.StartUDPServer("8081")

	// O QUE É ISTO????
	select {} // Bloqueio para manter o servidor rodando
}
