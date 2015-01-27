package main

import (
	"log"
	"net"
	"sieci/util"
	"strings"
	"sync"
)

type Server struct {
	workerListener *net.TCPListener
	clientListener *net.TCPListener
	Connections    []*util.Connection
}

func New(addr1 string, addr2 string) *Server {
	return &Server{
		workerListener: NewListener(addr1),
		clientListener: NewListener(addr2),
		Connections:    make([]*util.Connection, 0)}
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
		if err != nil {
			log.Fatal(err)
		}
		s.Connections = append(s.Connections, util.NewConnection(c))
	}
	quit <- true
}

func (s *Server) HandleClients(quit chan bool) {
	for {
		c, err := s.clientListener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go s.HandleRequest(util.NewConnection(c))
	}
	quit <- true
}

func (s *Server) HandleRequest(conn *util.Connection) {
	var mutex sync.Mutex
	handled := false

	for _, worker := range s.Connections {
		worker.Write("AVAILABLE")
		go func(worker *util.Connection) {
			cmd := worker.Read()

			if cmd[0] == "READY" {
				mutex.Lock()
				if !handled {
					handled = true
					mutex.Unlock()
					HandleUpload(conn, worker)
				} else {
					mutex.Unlock()
					worker.Write("ROLLBACK")
				}
			}
		}(worker)
	}
}

func HandleUpload(conn *util.Connection, worker *util.Connection) {
	line := conn.ReadLine()
	cmd := strings.Fields(line)
	if cmd[0] != "ENDSEQ" {
		log.Printf("FAIL\n")
		return
	}
	endseq := cmd[1]
	worker.WriteLine(line)
	line = conn.ReadLine()
	for strings.TrimSpace(line) != strings.TrimSpace(endseq) {
		worker.WriteLine(line)
		line = conn.ReadLine()
	}
	worker.WriteLine(line)
}
