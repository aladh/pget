package main

import (
	"github.com/ali-l/pget/download"
	"log"
	"os"
)

func main() {
	err := download.Run(os.Args[1])
	if err != nil {
		log.Fatalf("Download failed with error: %s", err)
	}
}
