package download

import (
	"errors"
	"fmt"
	"github.com/ali-l/pget/chunks"
	"log"
	"net/http"
	"os"
	"strings"
)

func Run(url string) error {
	res, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("error making HEAD request: %w", err)
	}

	if !supportsRangeRequests(res) {
		return errors.New("server does not support range requests")
	}

	contentLength := res.ContentLength

	filename := filename(url)

	log.Printf("Downloading %s (size %d)\n", filename, contentLength)

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	numParts := 2

	for i, chunk := range chunks.Build(url, contentLength, numParts, out) {
		err := chunk.Download()
		if err != nil {
			return fmt.Errorf("error downloading chunk %d: %w", i, err)
		}
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
