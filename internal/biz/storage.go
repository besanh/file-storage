package biz

import (
	"context"
	"time"
)

type StorageProvider interface {
	GetUploadUrl(ctx context.Context, key string, contentType string, expiry time.Duration) (string, error)
	GetDownloadUrl(ctx context.Context, key string, filename string, expiry time.Duration) (string, error)
	DeleteFile(ctx context.Context, key string) error
}
