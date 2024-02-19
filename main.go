package main

import (
	"fmt"
	"log"
	"net"
)




type NATRecord struct {
	Private string
	Public  string
}

var RecordCollection = []NATRecord{}

func main() {
	// Local address with a specific port
	localAddr, err := net.ResolveUDPAddr("udp", "192.168.1.71:8000")
	if err != nil {
		fmt.Println("Error resolving local address:", err)
		return
	}

	// Remote address
	remoteAddr, err := net.ResolveUDPAddr("udp", "20.193.146.188:5173")
	if err != nil {
		fmt.Println("Error resolving remote address:", err)
		return
	}

	// Dial UDP to establish a connection
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}

	// Send a message
	message := []byte(localAddr.String())
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error writing to UDP:", err)
		return
	}

	// Read a response
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error reading from UDP:", err)
		return
	}

	fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))
	conn.Close()

	connLis, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		// handle error
		log.Fatalln(err.Error())
	}
	defer connLis.Close()
	for {
		buf := make([]byte, 512)
		i, addr, err := connLis.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		RecordCollection = append(RecordCollection, NATRecord{
			Public:  addr.String(),
			Private: string(buf[0 : i-1]),
		})
		connLis.WriteTo([]byte("Replying to message!"), addr)
		log.Printf("%v %v", string(buf[0:i-1]), addr)
	}
}
