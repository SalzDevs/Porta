package main

import (
	"log"
	"net"
)

func main() {
	addr := "127.0.0.1:6432"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	log.Printf("porta listening on %s", addr)

	pool := NewPool(20)
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handleProxy(client, pool, "127.0.0.1:5432")
	}
}
