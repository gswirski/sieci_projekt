package main

import (
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
	fmt.Fprintf(conn, "OFFILESEQYOULLNEVERUSEITINYOURCODE\n")
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(conn, "OFFILESEQYOULLNEVERUSEITINYOURCODE\n")
	conn.Close()
}
