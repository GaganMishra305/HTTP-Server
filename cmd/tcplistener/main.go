package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		curline := ""
		for {
			data := make([]byte, 8)

			n, err := f.Read(data)
			if n > 0 {
				for _, ch := range string(data[:n]) {
					if ch == '\n' {
						lines <- curline
						curline = ""
					} else if ch != '\r' {
						curline += string(ch)
					}
				}
			}
			if err != nil {
				if curline != "" {
					lines <- curline
				}
				break
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
				if line != "" {
					fmt.Println(line)
				}
			}
			log.Println("Connection Closed")
		}(conn)

	}
}
