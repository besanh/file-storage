package biz

import (
	"context"
	db "file/internal/data/db/generated"

	"github.com/go-kratos/kratos/v2/log"
)

type PhysicalFileRepo interface {
	InsertPhysicalFile(ctx context.Context, file *db.InsertPhysicalFileParams) (*db.PhysicalFile, error)
	GetPhysicalFileByHash(ctx context.Context, fileHash string) (*db.PhysicalFile, error)
	UpdatePhysicalFileReferenceCount(ctx context.Context, file *db.UpdatePhysicalFileReferenceCountParams) (*db.PhysicalFile, error)
}

type PhysicalFileUsecase struct {
	physicalFileRepo PhysicalFileRepo
	authRepo         AuthRepo
	log              *log.Helper
}

func NewPhysicalFileUsecase(physicalFileRepo PhysicalFileRepo, authRepo AuthRepo, logger log.Logger) *PhysicalFileUsecase {
	return &PhysicalFileUsecase{
		physicalFileRepo: physicalFileRepo,
		authRepo:         authRepo,
		log:              log.NewHelper(logger),
	}
}
