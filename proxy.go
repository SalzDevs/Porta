package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
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

	pc := pool.Get(key)
	if pc != nil {
		log.Printf("[pool] reuse %s", key)
		if _, err := client.Write(pc.startup); err != nil {
			pc.conn.Close()
			return
		}
		defer pool.Put(key, pc)

		go func() {
			io.Copy(client, pc.conn)
			client.Close()
		}()

		if err := forwardClientMessages(buf, pc.conn); err != nil {
			log.Printf("client forward: %v", err)
		}
		return
	}

	log.Printf("[pool] dial %s", key)
	upstream, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		log.Printf("dial upstream: %v", err)
		return
	}

	if _, err := upstream.Write(lengthBuf[:]); err != nil {
		upstream.Close()
		return
	}
	if _, err := upstream.Write(payload); err != nil {
		upstream.Close()
		return
	}

	var startup bytes.Buffer
	if err := captureStartup(upstream, buf, client, &startup); err != nil {
		log.Printf("capture startup: %v", err)
		upstream.Close()
		return
	}

	go func() {
		io.Copy(client, upstream)
		client.Close()
	}()

	if err := forwardClientMessages(buf, upstream); err != nil {
		log.Printf("client forward: %v", err)
	}

	pool.Put(key, &PooledConn{conn: upstream, startup: startup.Bytes()})
}

func captureStartup(upstream net.Conn, clientBuf *bufio.Reader, client net.Conn, buf *bytes.Buffer) error {
	for {
		var msgType [1]byte
		if _, err := io.ReadFull(upstream, msgType[:]); err != nil {
			return err
		}

		var lengthBuf [4]byte
		if _, err := io.ReadFull(upstream, lengthBuf[:]); err != nil {
			return err
		}
		length := binary.BigEndian.Uint32(lengthBuf[:])

		payload := make([]byte, length-4)
		if _, err := io.ReadFull(upstream, payload); err != nil {
			return err
		}

		buf.Write(msgType[:])
		buf.Write(lengthBuf[:])
		buf.Write(payload)

		if _, err := client.Write(msgType[:]); err != nil {
			return err
		}
		if _, err := client.Write(lengthBuf[:]); err != nil {
			return err
		}
		if _, err := client.Write(payload); err != nil {
			return err
		}

		if msgType[0] == MsgAuthentication {
			authCode := binary.BigEndian.Uint32(payload)
			if authCode != AuthOK {
				pwType, err := clientBuf.ReadByte()
				if err != nil {
					return err
				}
				var pwLengthBuf [4]byte
				if _, err := io.ReadFull(clientBuf, pwLengthBuf[:]); err != nil {
					return err
				}
				pwLength := binary.BigEndian.Uint32(pwLengthBuf[:])
				pwPayload := make([]byte, pwLength-4)
				if _, err := io.ReadFull(clientBuf, pwPayload); err != nil {
					return err
				}

				if _, err := upstream.Write([]byte{pwType}); err != nil {
					return err
				}
				if _, err := upstream.Write(pwLengthBuf[:]); err != nil {
					return err
				}
				if _, err := upstream.Write(pwPayload); err != nil {
					return err
				}

				continue
			}
		}

		if msgType[0] == MsgReadyForQuery {
			return nil
		}
		if msgType[0] == MsgErrorResponse {
			return fmt.Errorf("startup failed")
		}
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

		if msgType == MsgTerminate {
			return nil
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
