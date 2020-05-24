package download

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	position := int64(0)

	for i := 0; i < numParts; i++ {
		part, err := downloadRange(url, position, position+contentLength/2-1)
		if err != nil {
			return fmt.Errorf("error downloading part: %w", err)
		}

		written, err := io.Copy(out, part)
		if err != nil {
			return fmt.Errorf("eror writing file: %w", err)
		}

		os.Remove(part.Name())

		position += contentLength / 2

		log.Printf("Wrote %d bytes\n", written)
	}

	return nil
}

func downloadRange(url string, start int64, end int64) (*os.File, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer res.Body.Close()

	out, err := ioutil.TempFile(os.TempDir(), "pget-")
	if err != nil {
		return nil, fmt.Errorf("error creating tempfile: %w", err)
	}

	written, err := io.Copy(out, res.Body)
	if err != nil {
		return nil, fmt.Errorf("error writing file: %w", err)
	}
	log.Printf("Downloaded %d chunk\n", written)

	_, err = out.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("error seeking to start of file: %w", err)
	}

	return out, nil
}

func supportsRangeRequests(res *http.Response) bool {
	acceptRanges := res.Header.Get("Accept-Ranges")
	return strings.Contains(acceptRanges, "bytes")
}

func filename(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}
