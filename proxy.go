package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
)

func handleProxy(client net.Conn, pool *Pool, upstreamAddr string) {
	defer client.Close()

	buf := bufio.NewReader(client)

	var lengthBuf [4]byte
	if _, err := io.ReadFull(buf, lengthBuf[:]); err != nil {
		return
	}
	length := binary.BigEndian.Uint32(lengthBuf[:])

	payload := make([]byte, length-4)
	if _, err := io.ReadFull(buf, payload); err != nil {
		return
	}

	fullMsg := append(lengthBuf[:], payload...)
	user, database, _, err := parse_startup(fullMsg)
	if err != nil {
		return
	}

	key := user + "/" + database
	log.Printf("[startup] user=%s database=%s key=%s", user, database, key)

	upstream, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		log.Printf("dial upstream: %v", err)
		return
	}
	defer upstream.Close()

	go func() {
		io.Copy(client, upstream)
		client.Close()
	}()

	if _, err := upstream.Write(lengthBuf[:]); err != nil {
		return
	}
	if _, err := upstream.Write(payload); err != nil {
		return
	}

	if err := forwardClientMessages(buf, upstream); err != nil {
		log.Printf("client forward: %v", err)
		return
	}
}

func forwardClientMessages(r *bufio.Reader, w io.Writer) error {
	for {
		msgType, err := r.ReadByte()
		if err != nil {
			return err
		}

		var lengthBuf [4]byte
		if _, err := io.ReadFull(r, lengthBuf[:]); err != nil {
			return err
		}
		length := binary.BigEndian.Uint32(lengthBuf[:])

		payload := make([]byte, length-4)
		if _, err := io.ReadFull(r, payload); err != nil {
			return err
		}

		if msgType == MsgQuery {
			fullMsg := append([]byte{msgType}, append(lengthBuf[:], payload...)...)
			if sql, err := parse_query(fullMsg); err == nil {
				log.Printf("[QUERY] %s", sql)
			}
		}

		if _, err := w.Write([]byte{msgType}); err != nil {
			return err
		}
		if _, err := w.Write(lengthBuf[:]); err != nil {
			return err
		}
		if _, err := w.Write(payload); err != nil {
			return err
		}

		if msgType == MsgTerminate {
			return nil
		}
	}
}
