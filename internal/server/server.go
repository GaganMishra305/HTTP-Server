package server

import (
    "bytes"
    "fmt"
    "io"
    "main/internal/request"
    "main/internal/response"
    "net"
    "sync"
)

type serverState string
type HandlerError struct {
    StatusCode 	response.StatusCode
    Message 	string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError 

const (
    StateInit    serverState = "init"
    StateServing serverState = "serving"
    StateClosed  serverState = "closed"
)

type Server struct {
    Port int
    handler Handler
    listener net.Listener
    state serverState
    mu sync.Mutex
}

func newServer(p int, handler Handler) (s *Server) {
    return &Server{
        Port:  p,
        handler: handler,  // Fix: actually assign the handler
        state: StateInit,
    }
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
    defer conn.Close()

    headers := response.GetDefaultHeaders(0)
    r, err := request.RequestFromReader(conn)
    if err != nil {
        response.WriteStatusLine(conn, response.StatusBadRequest)
        response.WriteHeaders(conn, headers)
        return
    }

    writer := bytes.NewBuffer([]byte{})
    handlerError := s.handler(writer, r)

    var body []byte = nil
    var status response.StatusCode = response.StatusOk
    if handlerError != nil {
        status = handlerError.StatusCode
        body = []byte(handlerError.Message)
    } else {
        body = writer.Bytes()
    }

    headers.Replace("Content-Length", fmt.Sprint(len(body)))

    response.WriteStatusLine(conn, status)
    response.WriteHeaders(conn, headers)
    conn.Write(body)
}

func runServer(s *Server, l net.Listener) {
    s.mu.Lock()
    s.state = StateServing
    s.mu.Unlock()
    
    for {
        conn, err := l.Accept()
        
        s.mu.Lock()
        closed := s.state == StateClosed
        s.mu.Unlock()
        
        if closed {
            return 
        }
        if err != nil {
            return
        }
        
        go runConnection(s, conn)
    }
}

func (s *Server) Close() {
    s.mu.Lock()
    s.state = StateClosed
    listener := s.listener
    s.mu.Unlock()
    
    if listener != nil {
        listener.Close()
    }
}

func Serve(port int, handler Handler) (*Server, error) {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return nil, err
    }
    s := newServer(port, handler)
    s.listener = listener
    go runServer(s, listener)

    return s, nil
}