package main

import (
	"fmt"
	"log"
	"net"

	"httpproto/internal/request"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Cannot open tcp connection, error: %s", err.Error())
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Cannot accept messages from TCP connection")
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Could not read HTTP request")
		}
		outputRequest(r)
	}
}

func outputRequest(r *request.Request) {
	fmt.Println("Request Line:")
	fmt.Printf("\tMethod: %s\n", r.RequestLine.Method)
	fmt.Printf("\tTarget: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("\tVersion: %s\n", r.RequestLine.HttpVersion)
}
