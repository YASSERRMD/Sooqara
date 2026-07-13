package store

import (
	"context"
	"io"
)

// Blob is a content-addressable storage interface.
type Blob interface {
	Put(ctx context.Context, key string, r io.Reader) (path string, err error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	URL(key string) string
}
