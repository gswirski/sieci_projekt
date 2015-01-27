package main

import (
	"log"
	"net"
	"os"
	"sieci/util"
	"time"
)

func HandleCommand(conn *util.Connection) {
	cmd := conn.Read()
	if cmd[0] == "ROLLBACK" {
		log.Printf("ROLLBACK\n")
		return
	}

	log.Printf("EXECUTE COMMAND\n")
	time.Sleep(5 * time.Second)
}

func HandleMaster(worker *util.Worker, conn *util.Connection) {
	for {
		_ = conn.Read()
		worker.Lock()

		conn.Write("READY")
		HandleCommand(conn)

		worker.Unlock()
	}
}

func main() {
	worker := util.Worker{}
	addrs := os.Args[1:]

	for _, addr := range addrs {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatal("w1", err)
		}

		go HandleMaster(&worker, util.NewConnection(c))
	}

	quit := make(chan bool)
	<-quit
}
