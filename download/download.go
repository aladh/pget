package download

import (
	"errors"
	"fmt"
	"github.com/ali-l/pget/chunks"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func Run(url string, numChunks int, verbose bool) error {
	startTime := time.Now()

	res, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("error making HEAD request: %w", err)
	}
	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing response body: %w", err)
	}

	if !supportsRangeRequests(res) {
		return errors.New("server does not support range requests")
	}

	contentLength := res.ContentLength

	filename := filename(url)

	if verbose {
		log.Printf("Downloading %s (%d bytes) in %d chunks\n", filename, contentLength, numChunks)
	}

	err = createFile(filename)
	if err != nil {
		return err
	}

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

	duration := time.Since(startTime).Seconds()

	if verbose {
		log.Printf("Finished in %f seconds. Average speed: %f MB/s\n", duration, float64(contentLength/1000000)/duration)
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

func supportsRangeRequests(res *http.Response) bool {
	acceptRanges := res.Header.Get("Accept-Ranges")
	return strings.Contains(acceptRanges, "bytes")
}

func filename(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}
