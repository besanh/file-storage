package data

import (
	"context"
	"database/sql"
	"fmt"

	"file/internal/biz"
	db "file/internal/data/db/generated"

	"github.com/go-kratos/kratos/v2/log"
)

type physicalFileRepo struct {
	data *Data
	log  *log.Helper
}

func NewPhysicalFileRepo(data *Data, logger log.Logger) biz.PhysicalFileRepo {
	return &physicalFileRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *physicalFileRepo) InsertPhysicalFile(ctx context.Context, file *db.InsertPhysicalFileParams) (*db.PhysicalFile, error) {
	resp, err := r.data.Queries(ctx).InsertPhysicalFile(ctx, *file)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file insert fail")
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
