package lm

import "context"

// RawSource fetches the upstream LM document (HTML snippet).
type RawSource interface {
	Fetch(ctx context.Context) ([]byte, error)
}
