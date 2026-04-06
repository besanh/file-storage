package biz

import (
	"context"
	"database/sql"
	"file/internal/common"
	db "file/internal/data/db/generated"
	"fmt"

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

func (uc *FileUsecase) CreateFile(ctx context.Context, input *CreateFileRequest) (CreateFileResponse, error) {
	ownerID, err := uuid.Parse(GetUserID(ctx))
	if err != nil {
		return CreateFileResponse{}, err
	}

	// Check permission
	if input.ParentID != nil {
		permCheck, err := uc.authRepo.CheckPermission(ctx, common.TypeUser.String(), GetUserID(ctx), common.RelationWriter.String(), common.TypeFolder.String(), input.ParentID.String())
		if err != nil {
			return CreateFileResponse{}, err
		}
		if !permCheck {
			return CreateFileResponse{}, fmt.Errorf("permission denied")
		}
	}

	var fileID uuid.UUID

	if err := uc.tm.ExecTx(ctx, func(ctx context.Context) error {
		// Physical file
		var psFile *db.PhysicalFile
		psFile, err = uc.physicalFileRepo.GetPhysicalFileByHash(ctx, input.FileHash)
		if err != nil {
			return err
		}
		if psFile != nil {
			psFile, err = uc.physicalFileRepo.UpdatePhysicalFileReferenceCount(ctx, &db.UpdatePhysicalFileReferenceCountParams{
				ID:             psFile.ID,
				ReferenceCount: sql.NullInt32{Int32: psFile.ReferenceCount.Int32 + 1, Valid: true},
			})
			if err != nil {
				return err
			}
		} else {
			psFile, err = uc.physicalFileRepo.InsertPhysicalFile(ctx, &db.InsertPhysicalFileParams{
				FileHash:       input.FileHash,
				SizeBytes:      input.FileSize,
				MimeType:       input.FileMimeType,
				StoragePath:    "",
				ReferenceCount: sql.NullInt32{Int32: 1, Valid: true},
			})
			if err != nil {
				return err
			}
		}

		// File
		var parentID uuid.NullUUID
		if input.ParentID != nil {
			parentID = uuid.NullUUID{UUID: *input.ParentID, Valid: true}
		}

		file := &db.InsertFileParams{
			OwnerID:             ownerID,
			ParentID:            parentID,
			Name:                input.Name,
			IsFolder:            input.IsFolder,
			PhysicalFileID:      uuid.NullUUID{UUID: psFile.ID, Valid: true},
			FileHash:            sql.NullString{String: input.FileHash, Valid: input.FileHash != ""},
			FileSize:            sql.NullInt64{Int64: input.FileSize, Valid: input.FileSize != 0},
			FileType:            sql.NullString{String: input.FileType, Valid: input.FileType != ""},
			FileExt:             sql.NullString{String: input.FileExt, Valid: input.FileExt != ""},
			FileMimeType:        sql.NullString{String: input.FileMimeType, Valid: input.FileMimeType != ""},
			FileVideoResolution: sql.NullString{String: input.FileVideoResolution, Valid: input.FileVideoResolution != ""},
			Status:              sql.NullString{String: input.Status, Valid: input.Status != ""},
		}
		fileNode, err := uc.fileRepo.InsertFile(ctx, file)
		if err != nil {
			return err
		}
		fileID = fileNode.ID

		return nil
	}); err != nil {
		return CreateFileResponse{}, err
	}

	return CreateFileResponse{
		ID: &fileID,
	}, nil
}
