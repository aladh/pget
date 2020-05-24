package download

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Run(url string) error {
	filename := filename(url)

	log.Printf("Downloading %s\n", filename)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Wrote %d bytes\n", written)

	return nil
}

func filename(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}
