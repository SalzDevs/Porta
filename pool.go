package main

import (
	"net"
	"sync"
)

type Pool struct {
	mu   sync.Mutex
	idle map[string][]net.Conn
}

func NewPool() *Pool {
	return &Pool{
		idle: make(map[string][]net.Conn),
	}
}

func (p *Pool) Get(key string) net.Conn {
	p.mu.Lock()
	defer p.mu.Unlock()

	conns := p.idle[key]
	if len(conns) == 0 {
		return nil
	}

	conn := conns[len(conns)-1]
	p.idle[key] = conns[:len(conns)-1]
	return conn
}

func (p *Pool) Put(key string, conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.idle[key] = append(p.idle[key], conn)
}
