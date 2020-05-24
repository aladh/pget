package main

import (
	"flag"
	"github.com/ali-l/pget/download"
	"log"
)

func main() {
	numChunks := flag.Int("c", 8, "number of chunks to download in parallel")
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()

	err := download.Run(flag.Arg(0), *numChunks, *verbose)
	if err != nil {
		log.Fatalf("Download failed with error: %s", err)
	}
}
