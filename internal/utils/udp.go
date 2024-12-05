package utils

import (
	"log"
	"net"
)

func WriteUDP(conn *net.UDPConn, udpAddr *net.UDPAddr, data []byte, infoMessage string, errorMessage string) {
	var err error

	if udpAddr == nil { // write UDP without using UDP address
		_, err = conn.Write(data)

	} else { // write UDP using UDP address
		_, err = conn.WriteToUDP(data, udpAddr)
	}

	if err != nil {
		log.Fatalln(Red+errorMessage, ":", err, Reset)
	}
	log.Println(Blue+infoMessage, Reset)
}

func ReadUDP(conn *net.UDPConn, infoMessage string, errorMessage string) (int, *net.UDPAddr, []byte) {
	newData := make([]byte, BUFFERSIZE)
	n, udpAddr, err := conn.ReadFromUDP(newData)
	if err != nil {
		log.Fatalln(Red+errorMessage, ":", err, Reset)
	}
	log.Println(Blue+infoMessage, Reset)
	return n, udpAddr, newData
}

func ResolveUDPAddrAndListen(ip string, port string) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", ip+":"+port)
	if err != nil {
		log.Fatalln(Red+"[ERROR 8] Unable to resolve address:", err, Reset)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln(Red+"[ERROR 59] Unable to initialize the agent:", err, Reset)
	}

	return conn
}

func ResolveUDPAddrAndDial(ip string, port string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+port)
	if err != nil {
		log.Fatalln(Red+"[ERROR 51] Unable to resolve address:", err, Reset)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalln(Red+"[ERROR 2] Unable to connect:", err, Reset)
	}
	return conn
}
