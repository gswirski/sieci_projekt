package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const timeout time.Duration = 5 * time.Second
const rate time.Duration = 1000 * time.Nanosecond
const maxConns int = 10

type Connection struct {
	server  *Server
	tcpConn *net.TCPConn
	busy    bool
	reader  *bufio.Reader
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
			&Connection{server: s, tcpConn: c, busy: false, reader: bufio.NewReader(c)})
		s.ConnMutex.Unlock()
	}

	quit <- true
}

func (s *Server) HandleRequest(conn *net.TCPConn) {
	var mutex sync.Mutex
	handled := false

	for _, worker := range s.Connections {
		fmt.Fprintf(worker.tcpConn, "AVAILABLE\n")
		go func(tcp net.Conn, reader *bufio.Reader) {
			log.Printf("worker PRE check %p\n", reader)
			l, err := reader.ReadString('\n')
			log.Printf("worker POST check\n")
			if err != nil {
				log.Print("[ERR] ", err)
				return
			}

			cmd := strings.Fields(l)
			if cmd[0] == "READY" {
				log.Printf("[read] READY\n")
				mutex.Lock()
				if !handled {
					handled = true
					mutex.Unlock()
					fmt.Fprintf(tcp, "UPLOAD dupa\n")
					log.Printf("[sent] UPLOAD dupa\n")
				} else {
					mutex.Unlock()
					fmt.Fprintf(tcp, "ROLLBACK\n")
					log.Printf("[sent] ROLLBACK\n")
				}
			}

		}(worker.tcpConn, worker.reader)
	}

	fmt.Fprintf(conn, "ERROR RESOURCE BUSY\n")
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
