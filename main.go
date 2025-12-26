package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser){
	curline := ""
	for {
		data := make([]byte, 8)

		n, err := f.Read(data)
		if err != nil {
			break
		}

		for _, ch := range string(data[:n]) {
			if ch == '\n' {
				fmt.Printf("Read: %s \n" , curline)
				curline = ""
			} else {
				curline += string(ch)
			}
		}
	}
}

func main() {
	f, err := os.Open("messages.txt")
	
	if err != nil {
		log.Fatal("File not there")
	}

	getLinesChannel(f)
}
