package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	url := "https://google.ca"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("google.html")
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
