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
			fmt.Printf("Request line:\n")
			fmt.Printf(" - Method: %s\n", req.RequestLine.Method)
			fmt.Printf(" - Target: %s\n", req.RequestLine.RequestTarget)
			fmt.Printf(" - Version: %s\n",req.RequestLine.HttpVersion)

			fmt.Println("Headers:")
			req.Headers.ForEach(func(n, v string) {
				fmt.Printf(" - %s: %s\n", n, v)
			}) 

			fmt.Println("Body:")
			fmt.Printf("%s\n", req.Body)
		}(conn)

	}
}
