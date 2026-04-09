package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID
	StoragePhotosUsed   int64
	StorageVideoUsed    int64
	StorageDocumentUsed int64
	StorageAudioUsed    int64
	StorageCompressUsed int64
	StorageOtherUsed    int64
	UpdatedAt           time.Time
}

func (u *User) TotalStorageUsed() int64 {
	return u.StoragePhotosUsed + u.StorageVideoUsed + u.StorageDocumentUsed +
		u.StorageAudioUsed + u.StorageCompressUsed + u.StorageOtherUsed
}

type UserRepo interface {
	GetUserStorage(ctx context.Context, id uuid.UUID) (*User, error)
}
