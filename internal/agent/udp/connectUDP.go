package udp

func ConnectUDP(serverAddr string) {
	conn := GetUDPConnection(serverAddr)

	defer conn.Close()

	HandleUDPConnection(conn)
}
