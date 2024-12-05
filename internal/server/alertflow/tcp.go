package alertflow

import (
	"log"
	utils "nms/internal/utils"
)

func StartTCPServer(port string) {
	listener := utils.ResolveTCPAddr("localhost", utils.SERVERTCP)

	for {

		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(utils.Red+"[ERROR 999] Unable to accept connection:", err, utils.Reset)
			continue
		}
		go handleTCPConnection(conn)

	}
}
