package main

import (
	"flag"
	"log"
)

func main() {
	client := flag.Bool("c", false, "client(c) or server(s)")
	server := flag.Bool("s", false, "client(c) or server(s)")
	port := flag.String("p", "1111", "port for local addr")
	relay := flag.String("rAddr", "192.168.1.71:5173", "relay address")

	flag.Parse()
	if *client {
		log.Println("Running client node")
		runNode(*port, *relay)
	} else if *server {
		log.Println("Running relay")
		runRelay(*port)
	} else {

		log.Fatal("-type is needed either c or s")
	}
}
