package server

import (
	"fmt"
	"io"
	"net"
)

type serverState string

const (
	StateInit    serverState = "init"
	StateServing serverState = "serving"
	StateClosed  serverState = "closed"
)

type Server struct {
	Port int

	state serverState
}

func newServer(p int) (s *Server) {
	return &Server{
		Port:  p,
		state: StateInit,
	}
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	body := "Hello World!"
	out := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
	conn.Write([]byte(out))
	conn.Close()
}

func runServer(s *Server, l net.Listener) {
	s.state = StateServing
	for {
		conn, err := l.Accept()
		if s.state == StateClosed {
			return 
		}
		if err != nil {
			return
		}
		
		go runConnection(s, conn)
	}
}

func (s *Server) Close() {
	s.state = StateClosed
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := newServer(port)
	go runServer(s, listener)

	return s, nil
}
