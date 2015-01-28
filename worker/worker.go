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

func HandleCommand(conn *util.Connection) error {
	cmd, err := conn.Read()
	if err != nil {
		return err
	}
	if cmd[0] == "ROLLBACK" {
		log.Printf("ROLLBACK\n")
		return nil
	}

	var result bytes.Buffer
	endseq := cmd[1]
	line, err := conn.ReadLine()
	if err != nil {
		return err
	}
	for strings.TrimSpace(line) != strings.TrimSpace(endseq) {
		result.WriteString(line)
		line, err = conn.ReadLine()
		if err != nil {
			return err
		}
	}

	conn.Write("RECEIVED")
	log.Printf(result.String())

	time.Sleep(1 * time.Second)

	conn.Write("ENDSEQ dupa")
	conn.Write("result")
	conn.Write("dupa")

	return nil
}

func HandleMaster(worker *util.Worker, conn *util.Connection, quit chan bool) {
	for {
		_, err := conn.Read()
		if err != nil {
			quit <- true
			return
		}
		worker.Lock()

		conn.Write("READY")
		err = HandleCommand(conn)
		if err != nil {
			quit <- true
			return
		}

		log.Printf("unlock\n")
		worker.Unlock()
	}

	quit <- true
}

func main() {
	worker := util.Worker{}
	addrs := os.Args[1:]

	quit := make(chan bool)

	for _, addr := range addrs {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatal("w1", err)
		}

		go HandleMaster(&worker, util.NewConnection(c), quit)
	}

	for _, _ = range addrs {
		<-quit
	}
}
