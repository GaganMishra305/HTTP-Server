package main

import (
	"fmt"
	"log"
	"main/internal/request"
	"net"
)

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
			req, err := request.RequestFromReader(c)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
			fmt.Println(req.Headers)
			log.Println("Connection Closed")
		}(conn)

	}
}
