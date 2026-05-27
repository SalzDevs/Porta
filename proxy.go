package main

import (
	"bufio"
	"encoding/binary"
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
		io.Copy(client, upstream)
		client.Close()
	}()

	buf := bufio.NewReader(client)

	if err := forwardStartup(buf, upstream); err != nil {
		log.Printf("startup forward: %v", err)
		return
	}

	if err := forwardClientMessages(buf, upstream); err != nil {
		log.Printf("client forward: %v", err)
		return
	}
}

func forwardStartup(r *bufio.Reader, w io.Writer) error {
	var lengthBuf [4]byte
	if _, err := io.ReadFull(r, lengthBuf[:]); err != nil {
		return err
	}
	length := binary.BigEndian.Uint32(lengthBuf[:])

	payload := make([]byte, length-4)
	if _, err := io.ReadFull(r, payload); err != nil {
		return err
	}

	if _, err := w.Write(lengthBuf[:]); err != nil {
		return err
	}
	if _, err := w.Write(payload); err != nil {
		return err
	}
	return nil
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
	}
}
