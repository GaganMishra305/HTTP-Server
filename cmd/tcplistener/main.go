package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string{
	lines := make(chan string)
	go func() {
		defer f.Close()
		curline := ""
		for {
			data := make([]byte, 8)

			n, err := f.Read(data)
			if err != nil {
				break
			}

			for _, ch := range string(data[:n]) {
				if ch == '\n' {
					lines <- curline
					curline = ""
				} else {
					curline += string(ch)
				}
			}
		}
		close(lines)
	}()

	return lines
}

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		log.Println("Connection Established")
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			for line := range getLinesChannel(c) {
				fmt.Println(line)
			}
			log.Println("Connection Closed")
		}(conn)

	}
}
