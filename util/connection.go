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

func (c *Connection) ReadLine() string {
	l, err := c.Reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[read] %s", l)
	return l
}

func (c *Connection) Read() []string {
	return strings.Fields(c.ReadLine())
}

func (c *Connection) WriteLine(cmd string) {
	log.Printf("[sent] %s", cmd)
	fmt.Fprintf(c.Conn, "%s", cmd)
}

func (c *Connection) Write(cmd string) {
	c.WriteLine(fmt.Sprintf("%s\n", cmd))
}