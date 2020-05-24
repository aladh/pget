package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filename(url))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("Wrote %d bytes\n", written)
}

func filename(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}
