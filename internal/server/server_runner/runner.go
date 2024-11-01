package main

import sc "nms/internal/server/server_config"

func main() {
	go sc.StartTCPServer("8080")
    sc.StartUDPServer("8081")

	select {} // Bloqueio para manter o servidor rodando
}
