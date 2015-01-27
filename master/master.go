package main

import "os"

func main() {
	s := New(os.Args[1], os.Args[2])
	defer s.Close()

	quit := make(chan bool)
	go s.HandleWorkers(quit)
	go s.HandleClients(quit)
	<-quit
}
