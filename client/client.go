package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":2001")
	if err != nil {
		log.Fatal(err)
	}

  fmt.Fprintf(conn, "Hello\n")
}
