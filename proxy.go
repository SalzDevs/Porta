package main

import (
	"io"
	"log"
	"net"
)

func handleProxy(client net.Conn, upstreamAddr string) {
	defer client.Close()

	upstream, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		log.Printf("dial upstream %s: %v", upstreamAddr, err)
		return
	}
	defer upstream.Close()

	go func() {
		io.Copy(upstream, client)
		if tcpConn, ok := upstream.(*net.TCPConn); ok {
			tcpConn.CloseWrite()
		}
	}()

	io.Copy(client, upstream)
	if tcpConn, ok := client.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
}
