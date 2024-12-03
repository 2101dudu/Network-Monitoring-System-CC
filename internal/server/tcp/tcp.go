package tcp

import (
	"log"
	utils "nms/internal/utils"
)

func StartTCPServer(port string) {
	listener := utils.ResolveTCPAddr("localhost", utils.SERVERTCP)

	for {

		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("[TCP] [ERROR] Unable to accept connection:", err)
			continue
		}
		go handleTCPConnection(conn)

	}
}
