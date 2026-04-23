package lmremote

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/praadit/lm-price/internal/domain/lm"
)

// RawSource fetches LM HTML from a configured HTTP URL.
type RawSource struct {
	URL    string
	Client *http.Client
}

var _ lm.RawSource = (*RawSource)(nil)

// NewRawSource returns an lm.RawSource backed by HTTP GET.
func NewRawSource(url string, timeout time.Duration) *RawSource {
	return &RawSource{
		URL: url,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Fetch implements lm.RawSource.
func (s *RawSource) Fetch(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
	if err != nil {
		return nil, err
	}

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
