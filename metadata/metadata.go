package metadata

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aladh/pget/config"
)

type Metadata struct {
	ContentLength         int64
	SupportsRangeRequests bool
}

func Fetch(url string) (*Metadata, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HEAD request: %w", err)
	}

	req.Header.Add("User-Agent", config.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing HEAD request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received bad response code %d", resp.StatusCode)
	}

	return &Metadata{
		ContentLength:         resp.ContentLength,
		SupportsRangeRequests: supportsRangeRequests(resp),
	}, nil
}

func supportsRangeRequests(resp *http.Response) bool {
	acceptRanges := resp.Header.Get("Accept-Ranges")
	return strings.Contains(acceptRanges, "bytes")
}
