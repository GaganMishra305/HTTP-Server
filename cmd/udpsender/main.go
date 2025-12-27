package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udp, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal("Cant establish UDP")
	}

	conn, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		log.Fatal("Cant establish UDP")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(line))
		if err != nil {
			log.Fatal(err)
		}
	}
}