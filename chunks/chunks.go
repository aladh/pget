package chunks

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aladh/pget/config"
)

type Chunk struct {
	url      string
	start    int64
	end      int64
	filename string
}

func Build(url string, contentLength int64, numChunks int, filename string) []Chunk {
	chunkSize := contentLength / int64(numChunks)

	position := int64(0)
	chunks := make([]Chunk, 0)

	for {
		if position+chunkSize > contentLength {
			chunks[len(chunks)-1].end += contentLength - position
			break
		}

		chunks = append(chunks, Chunk{
			url:      url,
			start:    position,
			end:      position + chunkSize,
			filename: filename,
		})

		position += chunkSize
	}

	return chunks
}

func (chunk *Chunk) Download() (err error) {
	req, err := http.NewRequest("GET", chunk.url, nil)
	if err != nil {
		return fmt.Errorf("error creating range request: %w", err)
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", chunk.start, chunk.end))
	req.Header.Add("User-Agent", config.UserAgent)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing range request: %w", err)
	}
	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing response body: %w", closeErr)
		}
	}()

	file, err := os.OpenFile(chunk.filename, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing file: %w", closeErr)
		}
	}()

	_, err = file.Seek(chunk.start, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking to start of chunk: %w", err)
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("error writing chunk: %w", err)
	}

	return nil
}
