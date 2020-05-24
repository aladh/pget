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
)

func Run(url string, numChunks int) error {
	res, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("error making HEAD request: %w", err)
	}
	defer res.Body.Close()

	if !supportsRangeRequests(res) {
		return errors.New("server does not support range requests")
	}

	contentLength := res.ContentLength

	filename := filename(url)

	log.Printf("Downloading %s (size %d) in %d chunks\n", filename, contentLength, numChunks)

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	err = out.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}

	wg := sync.WaitGroup{}

	for i, chunk := range chunks.Build(url, contentLength, numChunks, out) {
		wg.Add(1)
		go func(index int, chunk chunks.Chunk) {
			defer wg.Done()

			err := chunk.Download()
			if err != nil {
				log.Printf("error downloading chunk %d: %s\n", index, err)
			}
		}(i, chunk)
	}

	wg.Wait()

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
