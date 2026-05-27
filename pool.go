package main

import (
	"net"
	"sync"
	"time"
)

type PooledConn struct {
	conn     net.Conn
	startup  []byte
	lastUsed time.Time
}

type Pool struct {
	mu          sync.Mutex
	idle        map[string][]*PooledConn
	maxSize     int
	idleTimeout time.Duration
}

func NewPool(maxSize int, idleTimeout time.Duration) *Pool {
	p := &Pool{
		idle:        make(map[string][]*PooledConn),
		maxSize:     maxSize,
		idleTimeout: idleTimeout,
	}
	go p.sweep()
	return p
}

func (p *Pool) Get(key string) *PooledConn {
	p.mu.Lock()
	defer p.mu.Unlock()

	conns := p.idle[key]
	if len(conns) == 0 {
		return nil
	}

	pc := conns[len(conns)-1]
	p.idle[key] = conns[:len(conns)-1]
	return pc
}

func (p *Pool) Put(key string, pc *PooledConn) {
	pc.lastUsed = time.Now()

	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.idle[key]) >= p.maxSize {
		pc.conn.Close()
		return
	}
	p.idle[key] = append(p.idle[key], pc)
}

func (p *Pool) sweep() {
	for {
		time.Sleep(60 * time.Second)
		p.mu.Lock()
		for key, conns := range p.idle {
			alive := make([]*PooledConn, 0, len(conns))
			for _, pc := range conns {
				if time.Since(pc.lastUsed) > p.idleTimeout {
					pc.conn.Close()
				} else {
					alive = append(alive, pc)
				}
			}
			p.idle[key] = alive
		}
		p.mu.Unlock()
	}
}
