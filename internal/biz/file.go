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
	actorType, actorID, err := GetActorInfo(ctx)
	if err != nil {
		return CreateFileResponse{}, err
	}

	ownerID, err := uuid.Parse(actorID)
	if err != nil {
		if actorType == "service" {
			// For non-UUID client IDs, generate a deterministic UUID
			ownerID = uuid.NewMD5(uuid.NameSpaceOID, []byte(actorID))
		} else {
			return CreateFileResponse{}, fmt.Errorf("invalid user identity string: %v", err)
		}
	}

	// Strict Permission Check (Perform before any deep dive logic)
	// We only check permissions if we are inside a folder.
	// Root creation is currently allowed for all authenticated users to enable bootstrapping.
	if input.ParentID != nil {
		objectType := common.TypeFolder.String()
		objectID := input.ParentID.String()

		allowed, err := uc.authRepo.CheckPermission(
			ctx,
			actorType,
			actorID,
			common.PermissionCreate.String(),
			objectType,
			objectID,
		)
		if err != nil {
			return CreateFileResponse{}, fmt.Errorf("permission check failed: %w", err)
		}
		if !allowed {
			return CreateFileResponse{}, fmt.Errorf("permission denied: cannot create item in %s:%s", objectType, objectID)
		}
	}

	var fileID uuid.UUID
	var resType common.ResourceType

	if err := uc.tm.ExecTx(ctx, func(ctx context.Context) error {
		// File
		var parentID uuid.NullUUID
		if input.ParentID != nil {
			parentID = uuid.NullUUID{UUID: *input.ParentID, Valid: true}
		}

		file := db.InsertFileParams{
			OwnerID:  ownerID,
			ParentID: parentID,
			Name:     input.Name,
			IsFolder: input.IsFolder,
		}

		if !input.IsFolder {
			// Physical file handling ONLY for files
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
			file.PhysicalFileID = uuid.NullUUID{UUID: psFile.ID, Valid: true}

			// Metadata for non-folder nodes
			file.FileHash = sql.NullString{String: input.FileHash, Valid: input.FileHash != ""}
			file.FileSize = sql.NullInt64{Int64: input.FileSize, Valid: input.FileSize != 0}
			file.FileType = sql.NullString{String: input.FileType, Valid: input.FileType != ""}
			file.FileExt = sql.NullString{String: input.FileExt, Valid: input.FileExt != ""}
			file.FileMimeType = sql.NullString{String: input.FileMimeType, Valid: input.FileMimeType != ""}
			file.FileVideoResolution = sql.NullString{String: input.FileVideoResolution, Valid: input.FileVideoResolution != ""}
			file.Status = sql.NullString{String: input.Status, Valid: input.Status != ""}
		}

		fileNode, err := uc.fileRepo.InsertFile(ctx, &file)
		if err != nil {
			return err
		}
		fileID = fileNode.ID

		// Grant Permissions in SpiceDB
		resType = common.TypeFile
		if input.IsFolder {
			resType = common.TypeFolder
		}

		// Grant Ownership
		err = uc.authRepo.WriteRelationship(
			ctx,
			resType.String(),
			fileID.String(),
			common.RelationOwner.String(),
			actorType,
			actorID,
		)
		if err != nil {
			return fmt.Errorf("failed to grant ownership in SpiceDB: %w", err)
		}

		// Link to Parent (if applicable)
		if input.ParentID != nil {
			err = uc.authRepo.WriteRelationship(
				ctx,
				resType.String(),
				fileID.String(),
				common.RelationParent.String(),
				common.TypeFolder.String(),
				input.ParentID.String(),
			)
			if err != nil {
				return fmt.Errorf("failed to link to parent in SpiceDB: %w", err)
			}
		}

		return nil
	}); err != nil {
		// COMPENSATING TRANSACTION: If the DB transaction failed, cleanup SpiceDB
		if fileID != uuid.Nil {
			_ = uc.authRepo.DeleteRelationship(ctx, resType.String(), fileID.String(), common.RelationOwner.String(), actorType, actorID)
			if input.ParentID != nil {
				_ = uc.authRepo.DeleteRelationship(ctx, resType.String(), fileID.String(), common.RelationParent.String(), common.TypeFolder.String(), input.ParentID.String())
			}
		}
		return CreateFileResponse{}, err
	}

	return CreateFileResponse{
		ID: fileID,
	}, nil
}
