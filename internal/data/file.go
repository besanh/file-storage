package data

import (
	"context"
	"database/sql"
	"fmt"

	"file/internal/biz"
	db "file/internal/data/db/generated"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type fileRepo struct {
	data *Data
	log  *log.Helper
}

func NewFileRepo(data *Data, logger log.Logger) biz.FileRepo {
	return &fileRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *fileRepo) InsertFile(ctx context.Context, file *db.InsertFileParams) (*db.FileNode, error) {
	resp, err := r.data.Queries(ctx).InsertFile(ctx, *file)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file insert fail")
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *fileRepo) GetUserStorageUsed(ctx context.Context, ownerID uuid.UUID) (int64, error) {
	return r.data.Queries(ctx).GetUserStorageUsed(ctx, ownerID)
}
