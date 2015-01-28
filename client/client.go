package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	filename := os.Args[2]
	file, err := os.Open(strings.TrimSpace(filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintf(conn, "ENDSEQ OFFILESEQYOULLNEVERUSEITINYOURCODE\n")
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(conn, "OFFILESEQYOULLNEVERUSEITINYOURCODE\n")

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	cmd := strings.Fields(line)
	if cmd[0] != "RECEIVED" {
		log.Fatal("FAIL")
	}

	var result bytes.Buffer

	line, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	cmd = strings.Fields(line)
	if cmd[0] != "ENDSEQ" {
		log.Fatal("FAIL")
	}
	endseq := strings.TrimSpace(cmd[1])
	line, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	for strings.TrimSpace(line) != endseq {
		result.WriteString(line)
		line, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Print(result.String())
	conn.Close()
}
