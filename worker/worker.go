package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"sieci/util"
	"strings"
	"time"
)

func HandleCommand(conn *util.Connection) {
	cmd := conn.Read()
	if cmd[0] == "ROLLBACK" {
		log.Printf("ROLLBACK\n")
		return
	}

	var result bytes.Buffer
	endseq := cmd[1]
	line := conn.ReadLine()
	for strings.TrimSpace(line) != strings.TrimSpace(endseq) {
		result.WriteString(line)
		line = conn.ReadLine()
	}

	conn.Write("RECEIVED")
	log.Printf(result.String())

	time.Sleep(1 * time.Second)

	conn.Write("ENDSEQ dupa")
	conn.Write("result")
	conn.Write("dupa")
}

func HandleMaster(worker *util.Worker, conn *util.Connection) {
	for {
		_ = conn.Read()
		worker.Lock()

		conn.Write("READY")
		HandleCommand(conn)

		log.Printf("unlock\n")
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
