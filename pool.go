package main

import (
	"net"
	"sync"
)

type PooledConn struct {
	conn    net.Conn
	startup []byte
}

type Pool struct {
	mu   sync.Mutex
	idle map[string][]*PooledConn
}

func NewPool() *Pool {
	return &Pool{
		idle: make(map[string][]*PooledConn),
	}
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
	p.mu.Lock()
	defer p.mu.Unlock()

	p.idle[key] = append(p.idle[key], pc)
}
