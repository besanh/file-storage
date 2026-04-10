package data

import (
	"context"
	"file/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) GetUserStorage(ctx context.Context, id uuid.UUID) (*biz.User, error) {
	stats, err := r.data.Queries(ctx).GetUserStorageBreakdown(ctx, id)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:                  id,
		StoragePhotosUsed:   stats.Photos,
		StorageVideoUsed:    stats.Videos,
		StorageDocumentUsed: stats.Documents,
		StorageAudioUsed:    stats.Audio,
		StorageCompressUsed: stats.Compress,
		StorageOtherUsed:    stats.Other,
	}, nil
}
