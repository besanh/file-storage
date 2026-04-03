package biz

import (
	"context"
	db "file/internal/data/db/generated"

	"github.com/go-kratos/kratos/v2/log"
)

type PhysicalFileRepo interface {
	InsertPhysicalFile(ctx context.Context, file *db.InsertPhysicalFileParams) (*db.PhysicalFile, error)
}

type PhysicalFileUsecase struct {
	physicalFileRepo PhysicalFileRepo
	authRepo         AuthRepo
	tm               Transaction
	log              *log.Helper
}

func NewPhysicalFileUsecase(physicalFileRepo PhysicalFileRepo, authRepo AuthRepo, tm Transaction, logger log.Logger) *PhysicalFileUsecase {
	return &PhysicalFileUsecase{
		physicalFileRepo: physicalFileRepo,
		authRepo:         authRepo,
		tm:               tm,
		log:              log.NewHelper(logger),
	}
}
