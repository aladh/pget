package main

import (
	"flag"
	"github.com/ali-l/pget/download"
	"log"
)

func main() {
	numChunks := flag.Int("chunks", 8, "number of chunks to download in parallel")
	flag.Parse()

	err := download.Run(flag.Arg(0), *numChunks)
	if err != nil {
		log.Fatalf("Download failed with error: %s", err)
	}
}
