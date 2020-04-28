package storage

import (
	"context"
	"io"
)

// Storage interface
type Storage interface {
	Close() error
	Delete(ctx context.Context, objName string) error
	Write(ctx context.Context, r io.Reader, objName string, public bool) error
}
