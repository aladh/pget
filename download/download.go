package download

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ali-l/pget/chunks"
	"github.com/ali-l/pget/metadata"
)

func Run(url string, numChunks int, verbose bool) error {
	startTime := time.Now()

	meta, err := metadata.Fetch(url)
	if err != nil {
		return fmt.Errorf("error fetching metadata: %w", err)
	}

	if !meta.SupportsRangeRequests {
		return errors.New("server does not support range requests")
	}

	filename := filename(url)

	if verbose {
		log.Printf("Downloading %s (%d bytes) in %d chunks\n", filename, meta.ContentLength, numChunks)
	}

	err = createFile(filename)
	if err != nil {
		return err
	}

	downloadChunks(url, numChunks, verbose, meta.ContentLength, filename)

	duration := time.Since(startTime).Seconds()

	if verbose {
		log.Printf("Finished in %f seconds. Average speed: %f MB/s\n", duration, float64(meta.ContentLength/1000000)/duration)
	}

	return nil
}

func downloadChunks(url string, numChunks int, verbose bool, contentLength int64, filename string) {
	wg := sync.WaitGroup{}

	for i, chunk := range chunks.Build(url, contentLength, numChunks, filename) {
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

func filename(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}
