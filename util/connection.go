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

func (c *Connection) ReadLine() (string, error) {
	l, err := c.Reader.ReadString('\n')
	if err != nil {
		log.Print(err)
	}
	log.Printf("[read] %s", l)
	return l, err
}

func (c *Connection) Read() ([]string, error) {
	l, err := c.ReadLine()
	return strings.Fields(l), err
}

func (c *Connection) WriteLine(cmd string) {
	log.Printf("[sent] %s", cmd)
	fmt.Fprintf(c.Conn, "%s", cmd)
}

func (c *Connection) Write(cmd string) {
	c.WriteLine(fmt.Sprintf("%s\n", cmd))
}

func CopyData(src *Connection, dst *Connection) error {
	line, err := src.ReadLine()
	if err != nil {
		return err
	}
	cmd := strings.Fields(line)
	if cmd[0] != "ENDSEQ" {
		dst.Write("ERROR")
		return nil
	}
	endseq := cmd[1]
	dst.WriteLine(line)
	line, err = src.ReadLine()
	if err != nil {
		return err
	}
	for strings.TrimSpace(line) != strings.TrimSpace(endseq) {
		dst.WriteLine(line)
		line, err = src.ReadLine()
		if err != nil {
			return err
		}
	}
	dst.WriteLine(line)

	return nil
}
