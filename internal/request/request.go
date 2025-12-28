package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string
const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	state parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var SEPARATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func (rl *RequestLine) validHTTP() bool {
	return rl.HttpVersion == "1.1"
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {	
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP"{
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{	
		Method: string(parts[0]),
		RequestTarget:string(parts[1]),
		HttpVersion: string(httpParts[1]),
	}

	if !rl.validHTTP() {
		return nil, 0, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	return rl, read, nil
}

func (r* Request) parse(data []byte) (int, error) {
	read := 0

outer: 
	for{
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateDone

		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r* Request) done() bool{
	return r.state == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
 
	// buffer could get overloaded
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readn, err := request.parse(buf[:bufLen + n])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readn : bufLen])
		bufLen -= readn
	}

	return request, nil
}
