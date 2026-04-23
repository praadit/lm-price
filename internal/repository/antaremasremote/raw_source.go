package antaremasremote

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/praadit/lm-price/internal/domain/antaremas"
)

// RawSource fetches Antaremas HTML from a configured HTTP URL.
type RawSource struct {
	URL    string
	Client *http.Client
}

var _ antaremas.RawSource = (*RawSource)(nil)

// NewRawSource returns an antaremas.RawSource backed by HTTP GET.
func NewRawSource(url string, timeout time.Duration) *RawSource {
	return &RawSource{
		URL: url,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Fetch implements antaremas.RawSource.
func (s *RawSource) Fetch(ctx context.Context) ([]byte, error) {
	// Antaremas often protects the normal HTML page with a JS challenge.
	// Their WordPress JSON endpoint returns the page's HTML content without needing JS.
	if looksLikeHargaEmasPage(s.URL) {
		if html, err := s.fetchWPPageHTML(ctx); err == nil && len(html) > 0 {
			return html, nil
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
	if err != nil {
		return nil, err
	}
	// Antaremas blocks some default clients; send a browser-like UA.
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

func looksLikeHargaEmasPage(url string) bool {
	u := strings.ToLower(strings.TrimSpace(url))
	return strings.Contains(u, "antaremas.com") && strings.Contains(u, "/harga-emas")
}

type wpPage struct {
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
}

func (s *RawSource) fetchWPPageHTML(ctx context.Context) ([]byte, error) {
	// Known page id for /harga-emas/ as advertised by the Link header.
	const wpURL = "https://antaremas.com/wp-json/wp/v2/pages/112"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wpURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("fetch %s: status %d: %s", wpURL, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var p wpPage
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	if strings.TrimSpace(p.Content.Rendered) == "" {
		return nil, fmt.Errorf("fetch %s: empty rendered content", wpURL)
	}
	return []byte(p.Content.Rendered), nil
}
