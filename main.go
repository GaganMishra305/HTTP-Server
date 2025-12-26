package main

import (
	"fmt"
	"io"
	"log"
	"os"
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
	f, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal("File not there")
	}

	lines := getLinesChannel(f)
	for {
		val, ok := <-lines
		if !ok {
            break
        }
        fmt.Println("read:", val)
	}
}
