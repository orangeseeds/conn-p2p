package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/orangeseeds/holepunching/pkg/relay"
	"github.com/orangeseeds/holepunching/pkg/utils"
)

func main() {
	port := flag.String("port", ":8000", "Port Number")
	flag.Parse()

	// Local address with a specific port
	localAddr, err := net.ResolveUDPAddr("udp", "192.168.1.71"+*port)
	if err != nil {
		fmt.Println("Error resolving local address:", err)
		return
	}

	// Remote address
	remoteAddr, err := net.ResolveUDPAddr("udp", "192.168.1.71:8080")
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
	message, _ := utils.Encode(relay.Message{
		Src:  *localAddr,
		Dest: *remoteAddr,
		Type: relay.CONN,
	})
	_, err = conn.Write(message.Bytes())
	if err != nil {
		fmt.Println("Error writing to UDP:", err)
		return
	}

	records := make([]relay.NATRecord, 10)
	// Read a response
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error reading from UDP:", err)
		return
	}
    

	err = utils.Decode(buf[0:n], &records)
	if err != nil {
		log.Fatalln("Decode Error", err)
	}
	conn.Close()

	// connLis, err := net.ListenPacket("udp", *port)
	// if err != nil {
	// 	// handle error
	// 	log.Fatalln(err.Error())
	// }
	// defer connLis.Close()
	// for {
	// 	buf := make([]byte, 512)
	// 	i, addr, err := connLis.ReadFrom(buf)
	// 	if err != nil {
	// 		log.Println(err)
	// 		continue
	// 	}
	//
	// 	// connLis.WriteTo([]byte("Replying to message!"), addr)
	// 	log.Printf("%v %v", string(buf[0:i-1]), addr)
	// }
}
