package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
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

	idleTimeout := 10 * time.Minute
	if s := os.Getenv("PORTA_IDLE_TIMEOUT"); s != "" {
		if d, err := time.ParseDuration(s); err == nil && d > 0 {
			idleTimeout = d
		}
	}

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("shutting down...")
		cancel()
		listener.Close()
	}()

	pool := NewPool(ctx, poolSize, idleTimeout)
	var wg sync.WaitGroup

	log.Printf("porta listening on %s", listen)

	for {
		client, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Printf("accept: %v", err)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleProxy(ctx, client, pool, upstream)
		}()
	}

	wg.Wait()
	log.Println("shutdown complete")
}
