package util

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Worker struct {
	mutex sync.Mutex
}

func (w *Worker) Lock() {
	w.mutex.Lock()
}

func (w *Worker) Unlock() {
	w.mutex.Unlock()
}

type Connection struct {
	Conn   net.Conn
	Reader *bufio.Reader
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{Conn: conn, Reader: bufio.NewReader(conn)}
}

func (c *Connection) Read() []string {
	l, err := c.Reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	cmd := strings.Fields(l)
	log.Printf("[read] %s", cmd)
	return cmd
}

func (c *Connection) Write(cmd string) {
	log.Printf("[sent] %s\n", cmd)
	fmt.Fprintf(c.Conn, "%s\n", cmd)
}
