package main

import (
	"bytes"
	"log"
	"net"
)

type Message struct {
	Type uint
	Data []byte
}

func main() {
	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln("Error resolving laddr:", err)
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln("Error listening:", err)
	}
	log.Println("Listening for UDP packets on port 8080")
	for {
		data := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println("Read Error:", err)
			continue
		}
		log.Printf("Read %d bytes: %v", n, data[0:n])

		handler(*conn, data[0:n], addr)
	}
}

func handler(conn net.UDPConn, data []byte, addr *net.UDPAddr) {
	buf := new(bytes.Buffer)
	buf.Write([]byte("Received!"))
	conn.WriteToUDP([]byte{0x1}, addr)
	conn.WriteToUDP(buf.Bytes(), addr)
}
