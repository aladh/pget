package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Run(url string) error {
	contentLength, err := contentLength(url)
	if err != nil {
		return err
	}

	filename := filename(url)

	log.Printf("Downloading %s (size %d)\n", filename, contentLength)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("eror writing file: %w", err)
	}

	log.Printf("Wrote %d bytes\n", written)

	return nil
}

func filename(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}

func contentLength(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("error making HEAD request: %w", err)
	}

	return resp.ContentLength, nil
}
