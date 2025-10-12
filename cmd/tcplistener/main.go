package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
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
		fmt.Println("Connections are being accepted")
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	c := make(chan string)
	go func() {
		defer func(c chan string) {
			close(c)
			fmt.Println("Closing channel")
		}(c)
		defer f.Close()
		var currentLine string
		for {
			b := make([]byte, 8)
			_, err := f.Read(b)
			if err == io.EOF {
				c <- currentLine
				break
			}
			if n := bytes.Index(b, []byte("\n")); n > -1 {
				currentLine += string(b[:n])
				c <- currentLine
				currentLine = string(b[n+1:])
				continue
			}
			currentLine += string(b)
		}
	}()

	return c
}
