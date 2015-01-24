package main

func main() {
	s := New(":2000", ":2001")
	defer s.Close()

  quit := make(chan bool)
	go s.HandleWorkers(quit)
  go s.HandleClients(quit)
  <- quit
}
