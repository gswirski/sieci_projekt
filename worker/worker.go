package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"log"
	"net"
	"strings"
)

func readUntil(reader *bufio.Reader, endseq string) string {
	var result bytes.Buffer

	log.Printf("chuj\n")

	line, err := reader.ReadString('\n')

	log.Printf("[read] %s, %s\n", line, endseq)

	for line != endseq {
		log.Printf("[read] %s", line)

		if err != nil {
			log.Fatal(err)
		}

		result.WriteString(line)

		line, err = reader.ReadString('\n')
	}

	return result.String()
}

func main() {
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(conn)
	for {
		l, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[received] %s", l)

		cmd := strings.Fields(l)
		log.Printf("fields: %q\n", cmd)

		var response string

		if cmd[0] == "AVAILABLE" {
			response = "OK\n"
		} else if cmd[0] == "ENDSEQ" {
			code := readUntil(reader, cmd[1])
			response = "OK\n"
			log.Printf("[loaded] %s", code)
		} else if cmd[0] == "SHUTDOWN" {
			response = "SHUTTING DOWN\n"
		} else {
			response = "ERROR\n"
		}

		fmt.Fprintf(conn, "OK\n")
		log.Printf("[sent] %s", response)
	}
}
