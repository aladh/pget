package metadata

import (
	"fmt"
	"net/http"
	"strings"
)

type Metadata struct {
	ContentLength         int64
	Filename              string
	SupportsRangeRequests bool
}

func New(url string) (*Metadata, error) {
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
		Filename:              filename(url),
		SupportsRangeRequests: supportsRangeRequests(res),
	}, nil
}

func supportsRangeRequests(res *http.Response) bool {
	acceptRanges := res.Header.Get("Accept-Ranges")
	return strings.Contains(acceptRanges, "bytes")
}

func filename(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}
