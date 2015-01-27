package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const timeout time.Duration = 5 * time.Second
const rate time.Duration = 1000 * time.Nanosecond
const maxConns int = 10

type Connection struct {
	server    *Server
	tcpConn   *net.TCPConn
	busy      bool
}

func (c *Connection) Handle() {

}

type Server struct {
	workerListener *net.TCPListener
	clientListener *net.TCPListener
	ConnMutex      sync.Mutex
	Connections    []*Connection
}

func New(addr1 string, addr2 string) *Server {
	return &Server{
		workerListener: NewListener(addr1),
		clientListener: NewListener(addr2),
		Connections:    make([]*Connection, 0)}
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
		if len(s.Connections) >= maxConns {
			c.Close()
			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		s.ConnMutex.Lock()
		s.Connections = append(s.Connections,
			&Connection{server: s, tcpConn: c, busy: false})
		s.ConnMutex.Unlock()
	}

	quit <- true
}

func (s *Server) HandleRequest(conn *net.TCPConn) {
	for _, worker := range s.Connections {
		if !worker.busy {
			worker.busy = true

			reader := bufio.NewReader(conn)
			for {
				r, err := reader.ReadString('\n')
				if err != nil {
					log.Printf("error: %s\n", err)
					break
				}
				log.Printf("passing %s", r)

				fmt.Fprintf(worker.tcpConn, r)
			}

			worker.busy = true
			return
		}
	}

	fmt.Fprintf(conn, "ERROR RESOURCE BUSY\n")
	conn.Close()
}

func (s *Server) HandleClients(quit chan bool) {
	for {
		c, err := s.clientListener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}

		go s.HandleRequest(c)
	}

	quit <- true
}
