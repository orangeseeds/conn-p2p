package main

import (
	"flag"
	"log"
)

func main() {
	client := flag.Bool("c", false, "node as client(c)")
	server := flag.Bool("s", false, "node as server(s)")
    port := flag.String("p", ":1111", "port for local addr")
	relay := flag.String("rAddr", ":5173", "relay address")

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
