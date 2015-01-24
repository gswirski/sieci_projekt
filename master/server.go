package main

import (
	"log"
	"net"
	"sync"
	"time"
)

const timeout time.Duration = 5 * time.Second
const rate time.Duration = 1000 * time.Nanosecond
const maxConns int = 10

type Connection struct {
	server  *Server
	tcpConn *net.TCPConn
}

func NewConnection(s *Server, c *net.TCPConn) *Connection {
	s.ConnMutex.Lock()
	s.Connections += 1
	s.ConnMutex.Unlock()
	return &Connection{server: s, tcpConn: c}
}

func (c *Connection) Handle() {

}

type Server struct {
	workerListener *net.TCPListener
	clientListener *net.TCPListener
	ConnMutex      sync.Mutex
	Connections    int
	ConnectionPool map[string]*Connection
}

func New(addr1 string, addr2 string) *Server {

	return &Server{
    workerListener: NewListener(addr1),
    clientListener: NewListener(addr2),
		Connections:    0,
		ConnectionPool: make(map[string]*Connection)}
}

func NewListener(addr string) *net.TCPListener {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

  return l.(*net.TCPListener)
}

func (s *Server) Close() {
	s.workerListener.Close()
	s.clientListener.Close()
}

func (s *Server) HandleWorkers(quit chan bool) {
	for {
		c, err := s.workerListener.AcceptTCP()
		if s.Connections >= maxConns {
			c.Close()
			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		NewConnection(s, c)
	}

  quit <- true
}

func (s *Server) HandleClients(quit chan bool) {

}
