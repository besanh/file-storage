package biz

import (
	"context"
	"database/sql"
	db "file/internal/data/db/generated"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type FileRepo interface {
	InsertFile(ctx context.Context, file *db.InsertFileParams) (*db.FileNode, error)
}

type FileUsecase struct {
	fileRepo         FileRepo
	physicalFileRepo PhysicalFileRepo
	authRepo         AuthRepo
	tm               Transaction
	log              *log.Helper
}

func NewFileUsecase(fileRepo FileRepo, physicalFileRepo PhysicalFileRepo, authRepo AuthRepo, tm Transaction, logger log.Logger) *FileUsecase {
	return &FileUsecase{
		fileRepo:         fileRepo,
		physicalFileRepo: physicalFileRepo,
		authRepo:         authRepo,
		tm:               tm,
		log:              log.NewHelper(logger),
	}
}

func (uc *FileUsecase) CreateFile(ctx context.Context, parentID *uuid.UUID, name string, isFolder bool, fileHash string, fileSize int64, fileType string, fileExt string, fileMimeType string, fileVideoResolution string, status string) error {
	ownerID, err := uuid.Parse(GetUserID(ctx))
	if err != nil {
		return err
	}

	if err := uc.tm.ExecTx(ctx, func(ctx context.Context) error {

		file := &db.InsertFileParams{
			OwnerID:  ownerID,
			ParentID: uuid.NullUUID{UUID: *parentID, Valid: parentID != nil},
			// PhysicalFileID:      uuid.NullUUID{UUID: physicalFileID, Valid: true},
			Name:                name,
			IsFolder:            isFolder,
			FileHash:            sql.NullString{String: fileHash, Valid: fileHash != ""},
			FileSize:            sql.NullInt64{Int64: fileSize, Valid: fileSize != 0},
			FileType:            sql.NullString{String: fileType, Valid: fileType != ""},
			FileExt:             sql.NullString{String: fileExt, Valid: fileExt != ""},
			FileMimeType:        sql.NullString{String: fileMimeType, Valid: fileMimeType != ""},
			FileVideoResolution: sql.NullString{String: fileVideoResolution, Valid: fileVideoResolution != ""},
			Status:              sql.NullString{String: status, Valid: status != ""},
		}
		_, err := uc.fileRepo.InsertFile(ctx, file)
		if err != nil {
			return err
		}

		physicalFile := &db.InsertPhysicalFileParams{
			FileHash:       fileHash,
			SizeBytes:      fileSize,
			MimeType:       fileMimeType,
			StoragePath:    "",
			ReferenceCount: sql.NullInt32{Int32: 1, Valid: true},
		}
		_, err = uc.physicalFileRepo.InsertPhysicalFile(ctx, physicalFile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
