package main

import (
	"fmt"
	"os"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(conn, "Hello\n")
}
