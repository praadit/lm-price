package galeri24remote

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/praadit/lm-price/internal/domain/galeri24"
)

// RawSource fetches Galeri24 HTML from a configured HTTP URL.
type RawSource struct {
	URL    string
	Client *http.Client
}

var _ galeri24.RawSource = (*RawSource)(nil)

// NewRawSource returns a galeri24.RawSource backed by HTTP GET.
func NewRawSource(url string, timeout time.Duration) *RawSource {
	return &RawSource{
		URL: url,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Fetch implements galeri24.RawSource.
func (s *RawSource) Fetch(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("fetch %s: status %d: %s", s.URL, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return io.ReadAll(resp.Body)
}
