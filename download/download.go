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

type Download struct {
	url       string
	numChunks int
	filename  string
	verbose   bool
	waitGroup sync.WaitGroup
}

func New(url string, numChunks int, verbose bool) *Download {
	return &Download{
		url:       url,
		numChunks: numChunks,
		filename:  filename(url),
		verbose:   verbose,
		waitGroup: sync.WaitGroup{},
	}
}

func (d *Download) Run() error {
	startTime := time.Now()

	meta, err := metadata.Fetch(d.url)
	if err != nil {
		return fmt.Errorf("error fetching metadata: %w", err)
	}

	if !meta.SupportsRangeRequests {
		return errors.New("server does not support range requests")
	}

	if d.verbose {
		log.Printf("Downloading %s (%d bytes) in %d chunks\n", d.filename, meta.ContentLength, d.numChunks)
	}

	err = d.createFile()
	if err != nil {
		return err
	}

	d.downloadChunks(meta.ContentLength)
	d.waitGroup.Wait()

	duration := time.Since(startTime).Seconds()

	if d.verbose {
		log.Printf("Finished in %f seconds. Average speed: %f MB/s\n", duration, float64(meta.ContentLength/1000000)/duration)
	}

	return nil
}

func (d *Download) downloadChunks(contentLength int64) {
	for i, chunk := range chunks.Build(d.url, contentLength, d.numChunks, d.filename) {
		d.waitGroup.Add(1)

		go func(index int, chunk chunks.Chunk) {
			defer d.waitGroup.Done()

			err := chunk.Download()
			if err != nil {
				log.Printf("error downloading chunk %d: %s\n", index, err)
			}

			if d.verbose {
				log.Printf("Downloaded chunk %d\n", index)
			}
		}(i, chunk)
	}
}

func (d *Download) createFile() error {
	file, err := os.Create(d.filename)
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
