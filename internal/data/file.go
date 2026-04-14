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

func (r *fileRepo) GetRecentFiles(ctx context.Context, ownerID uuid.UUID, limit int32) ([]*db.FileNode, error) {
	nodes, err := r.data.Queries(ctx).GetRecentFiles(ctx, db.GetRecentFilesParams{
		OwnerID: ownerID,
		Limit:   limit,
	})
	if err != nil {
		return nil, err
	}

	res := make([]*db.FileNode, 0, len(nodes))
	for _, n := range nodes {
		// Use a local copy to avoid pointing to the same loop variable
		node := n
		res = append(res, &node)
	}
	return res, nil
}

func (r *fileRepo) GetFile(ctx context.Context, id uuid.UUID) (*db.FileNode, error) {
	resp, err := r.data.Queries(ctx).GetFile(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *fileRepo) GetUserFileCount(ctx context.Context, ownerID uuid.UUID) (int64, error) {
	return r.data.Queries(ctx).GetUserFileCount(ctx, ownerID)
}

func (r *fileRepo) GetUserFolderCount(ctx context.Context, ownerID uuid.UUID) (int64, error) {
	return r.data.Queries(ctx).GetUserFolderCount(ctx, ownerID)
}

func (r *fileRepo) GetRecentFilesReport(ctx context.Context, ownerID uuid.UUID) (int64, int64, error) {
	row, err := r.data.Queries(ctx).GetRecentFilesReport(ctx, ownerID)
	if err != nil {
		return 0, 0, err
	}
	return row.NewFilesCount, row.NewStorageUsed, nil
}
