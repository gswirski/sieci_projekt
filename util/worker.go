package main

import "sieci/util"

type Worker struct {
	*Server
}

func NewWorker(wAddr string, mAddr string) *Worker {
	return &Worker{
		Server: NewServer(wAddr),
		Master: &Connection{conn: net.Dial("tcp", mAddr).(*net.TCPConn)},
	}
}

func main() {
	worker := NewWorker(os.Args[1], os.Args[2])

}
