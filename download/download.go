package download

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ali-l/pget/chunks"
	"github.com/ali-l/pget/metadata"
)

func Run(url string, numChunks int, verbose bool) error {
	startTime := time.Now()

	meta, err := metadata.New(url)
	if err != nil {
		return fmt.Errorf("error finding metadata: %w", err)
	}

	if !meta.SupportsRangeRequests {
		return errors.New("server does not support range requests")
	}

	if verbose {
		log.Printf("Downloading %s (%d bytes) in %d chunks\n", meta.Filename, meta.ContentLength, numChunks)
	}

	err = createFile(meta.Filename)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for i, chunk := range chunks.Build(url, meta.ContentLength, numChunks, meta.Filename) {
		wg.Add(1)
		go func(index int, chunk chunks.Chunk) {
			defer wg.Done()

			err := chunk.Download()
			if err != nil {
				log.Printf("error downloading chunk %d: %s\n", index, err)
			}

			if verbose {
				log.Printf("Downloaded chunk %d\n", index)
			}
		}(i, chunk)
	}

	wg.Wait()

	duration := time.Since(startTime).Seconds()

	if verbose {
		log.Printf("Finished in %f seconds. Average speed: %f MB/s\n", duration, float64(meta.ContentLength/1000000)/duration)
	}

	return nil
}

func createFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}

	return nil
}
