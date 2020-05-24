package chunks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Chunk struct {
	url   string
	start int64
	end   int64
	out   *os.File
}

func Build(url string, contentLength int64, numChunks int, out *os.File) []Chunk {
	chunkSize := contentLength / int64(numChunks)

	position := int64(0)
	chunks := make([]Chunk, 0)

	for {
		if position+chunkSize > contentLength {
			chunks[len(chunks)-1].end += contentLength - position
			break
		}

		file, err := os.OpenFile(out.Name(), os.O_RDWR, 0666)
		if err != nil {
			log.Println("error opening file")
		}

		chunks = append(chunks, Chunk{
			url:   url,
			start: position,
			end:   position + chunkSize,
			out:   file,
		})

		position += chunkSize
	}

	return chunks
}

func (chunk *Chunk) Download() error {
	req, err := http.NewRequest("GET", chunk.url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", chunk.start, chunk.end))
	req.Close = true

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request: %w", err)
	}
	defer res.Body.Close()

	_, err = chunk.out.Seek(chunk.start, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking to start of file: %w", err)
	}

	_, err = io.Copy(chunk.out, res.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	err = chunk.out.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}

	log.Printf("Downloaded chunk (%d-%d)\n", chunk.start, chunk.end)

	return nil
}
