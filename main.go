package main

import (
	"flag"
	"log"

	"github.com/aladh/pget/download"
)

func main() {
	numChunks := flag.Int("c", 8, "number of chunks to download in parallel")
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()

	err := download.New(flag.Arg(0), *numChunks, *verbose).Run()
	if err != nil {
		log.Fatalf("Download failed with error: %s", err)
	}
}
