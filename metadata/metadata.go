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
	res, err := http.Head(url)
	if err != nil {
		return nil, fmt.Errorf("error making HEAD request: %w", err)
	}
	err = res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing response body: %w", err)
	}

	return &Metadata{
		ContentLength:         res.ContentLength,
		SupportsRangeRequests: supportsRangeRequests(res),
	}, nil
}

func supportsRangeRequests(res *http.Response) bool {
	acceptRanges := res.Header.Get("Accept-Ranges")
	return strings.Contains(acceptRanges, "bytes")
}
