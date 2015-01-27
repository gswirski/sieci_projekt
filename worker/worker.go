package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func readUntil(conn io.Reader, endseq string) string {
	var result bytes.Buffer

	for line, err := bufio.NewReader(conn).ReadString('\n'); line != endseq; {
		if err != nil {
			log.Fatal(err)
		}

		result.WriteString(line)
	}

	return result.String()
}

func main() {
	conn, err := net.Dial("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		cmd, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[received] %s", cmd)

		var response string

		if cmd == "AVAILABLE\n" {
			response = "OK\n"
		} else if cmd == "ENDSEQ" {
			seq := strings.Split(cmd, " ")[1]
			code := readUntil(conn, seq)
			response = "OK\n"
			log.Printf("[loaded] %s", code)
		} else if cmd == "SHUTDOWN\n" {
			response = "SHUTTING DOWN\n"
		} else {
			response = "ERROR\n"
		}

		fmt.Fprintf(conn, "OK\n")
		log.Printf("[sent] %s", response)
	}
}
