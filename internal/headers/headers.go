package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var rn = []byte("\r\n")

type Headers struct { 
	headers map[string]string
}

func isToken(str []byte) bool{
	for _, ch := range str {
		found := false
		if 	ch >= 'A' && ch <= 'Z' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= '0' && ch <= '9'{
			 	found = true 
			}

		switch ch {
			case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
				found = true
		}

		if !found {
			return false
		}
	}
	return true
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)

	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s, %s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

func NewHeaders() *Headers {
	return &Headers{
		headers : map[string]string{},
	} 
}

func parseHeader(fieldline []byte) (string, string, error) {
	parts := bytes.SplitN(fieldline, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Malformed filed line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("Malformed field name")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) ( int, bool,  error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		//empty header
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}

		name, value, err := parseHeader(data[read:idx+read])
		if err != nil {
			return 0, false, err
		}
		if !isToken([]byte(name)){
			return 0, false, fmt.Errorf("Malofrmed header name")
		}

		read += idx + len(rn)
		h.Set(name, value)
	}

	return read, done, nil
}