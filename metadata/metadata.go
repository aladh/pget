package metadata

import (
	"fmt"
	"net/http"
	"strings"
)

type Metadata struct {
	ContentLength         int64
	SupportsRangeRequests bool
}

func Fetch(url string) (*Metadata, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, fmt.Errorf("error making HEAD request: %w", err)
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
