package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	dst, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("Cannot resolve UDP address")
	}
	conn, err := net.DialUDP("udp", nil, dst)
	if err != nil {
		log.Fatal("Cannot stablish UDP connection")
	}
	b := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		i, err := b.ReadString(byte('\n'))
		if err != nil {
			log.Fatal("Cannot read from stdin")
		}
		_, err = conn.Write([]byte(i))
		if err != nil {
			log.Fatal("Cannot write to UDP connection")
		}
	}
}
