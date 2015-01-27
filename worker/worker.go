package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Worker struct {
	mutex sync.Mutex
}

type Connection struct {
	reader *bufio.Reader
	conn   net.Conn
}

func (w *Worker) Lock() {
	w.mutex.Lock()
}

func (w *Worker) Unlock() {
	w.mutex.Unlock()
}

func (c *Connection) Read() []string {
	l, err := c.reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	cmd := strings.Fields(l)
	log.Printf("[read] %s", cmd)
	return cmd
}

func (c *Connection) Write(cmd string) {
	log.Printf("[sent] %s\n", cmd)
	fmt.Fprintf(c.conn, "%s\n", cmd)
}

func HandleCommand(conn *Connection) {
	cmd := conn.Read()
	if cmd[0] == "ROLLBACK" {
		log.Printf("rollback\n")
		return
	}

	log.Printf("start reading...")
	log.Printf(" sleep...")
	time.Sleep(5 * time.Second)
	log.Printf(" wake up\n")

}

func HandleMaster(worker *Worker, conn *Connection) {
	for {
		_ = conn.Read()
		worker.Lock()

		conn.Write("READY")
		HandleCommand(conn)

		worker.Unlock()
	}
}

func main() {
	worker := Worker{}
	addrs := os.Args[1:]

	for _, addr := range addrs {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatal("w1", err)
		}

		go HandleMaster(&worker, &Connection{conn: c, reader: bufio.NewReader(c)})
	}

	quit := make(chan bool)
	<-quit
}
