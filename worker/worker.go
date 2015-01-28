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
	cmd, err := conn.Read()
	if err != nil {
		return
	}
	if cmd[0] == "ROLLBACK" {
		log.Printf("ROLLBACK\n")
		return
	}

	var result bytes.Buffer
	endseq := cmd[1]
	line, err := conn.ReadLine()
	if err != nil {
		return
	}
	for strings.TrimSpace(line) != strings.TrimSpace(endseq) {
		result.WriteString(line)
		line, err = conn.ReadLine()
		if err != nil {
			return
		}
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
		_, err := conn.Read()
		if err != nil {
			return
		}
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
