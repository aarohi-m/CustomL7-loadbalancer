package main

import (
	"log"
	"net"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func NewBackend(u *url.URL) *Backend {
	return &Backend{URL: u, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(u)}
}

type ServerPool struct {
	backends []*Backend
	current  uint64
}

func (s *ServerPool) AddBackend(b *Backend) { s.backends = append(s.backends, b) }

func (s *ServerPool) GetNext() *Backend {
	idx := atomic.AddUint64(&s.current, 1) % uint64(len(s.backends))
	return s.backends[idx]
}

func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		conn, err := net.DialTimeout("tcp", b.URL.Host, 2*time.Second)
		b.mux.Lock()
		b.Alive = (err == nil)
		if err == nil {
			conn.Close()
		}
		b.mux.Unlock()
		log.Printf("Check: %s Alive: %v", b.URL, b.Alive)
	}
}
