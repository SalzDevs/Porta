package main

import (
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	listen := os.Getenv("PORTA_LISTEN")
	if listen == "" {
		listen = "127.0.0.1:6432"
	}

	upstream := os.Getenv("PORTA_UPSTREAM")
	if upstream == "" {
		upstream = "127.0.0.1:5432"
	}

	poolSize := 20
	if s := os.Getenv("PORTA_POOL_SIZE"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			poolSize = n
		}
	}

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	log.Printf("porta listening on %s", listen)

	pool := NewPool(poolSize)
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handleProxy(client, pool, upstream)
	}
}
