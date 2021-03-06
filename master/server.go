package main

import (
	"log"
	"net"
	"sieci/util"
	"sync"
)

type Server struct {
	workerListener *net.TCPListener
	clientListener *net.TCPListener
	Connections    map[*util.Connection]bool
}

func New(addr1 string, addr2 string) *Server {
	return &Server{
		workerListener: NewListener(addr1),
		clientListener: NewListener(addr2),
		Connections:    make(map[*util.Connection]bool)}
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
		s.Connections[util.NewConnection(c)] = false
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

	for worker, busy := range s.Connections {
		if busy {
			continue
		}

		s.Connections[worker] = true
		worker.Write("AVAILABLE")
		go func(worker *util.Connection) {
			cmd, err := worker.Read()
			if err != nil {
				delete(s.Connections, worker)
				return
			}

			if cmd[0] == "READY" {
				mutex.Lock()
				if !handled {
					handled = true
					mutex.Unlock()
					err := HandleUpload(conn, worker)
					if err != nil {
						delete(s.Connections, worker)
						return
					}
				} else {
					mutex.Unlock()
					worker.Write("ROLLBACK")
				}
				s.Connections[worker] = false // free
			}
		}(worker)
	}
}

func HandleUpload(conn *util.Connection, worker *util.Connection) error {
	err := util.CopyData(conn, worker)
	if err != nil {
		return nil
	}

	cmd, err := worker.Read()
	if err != nil {
		return err
	}
	if cmd[0] != "RECEIVED" {
		conn.Write("ERROR")
		return nil
	}

	conn.Write("RECEIVED")
	err = util.CopyData(worker, conn)
	if err != nil {
		return err
	}

	return nil
}
