package service

import (
	"context"

	pb "file/api/file/v1"
	"file/internal/biz"

	"github.com/google/uuid"
)

type FileService struct {
	pb.UnimplementedFileServiceServer
	uc *biz.FileUsecase
}

func NewFileService(uc *biz.FileUsecase) *FileService {
	return &FileService{uc: uc}
}

func (s *FileService) CreateFile(ctx context.Context, req *pb.CreateFileRequest) (*pb.CreateFileResponse, error) {
	var parentIDPtr *uuid.UUID
	if req.ParentId != "" {
		parsedID, err := uuid.Parse(req.ParentId)
		if err != nil {
			return nil, err
		}
		parentIDPtr = &parsedID
	}

	resp, err := s.uc.CreateFile(ctx, &biz.CreateFileRequest{
		ParentID:            parentIDPtr,
		Name:                req.Name,
		IsFolder:            req.IsFolder,
		FileHash:            req.FileHash,
		FileSize:            req.FileSize,
		FileType:            req.FileType,
		FileExt:             req.FileExt,
		FileMimeType:        req.FileMimeType,
		FileVideoResolution: req.FileVideoResolution,
		Status:              req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateFileResponse{
		Id: resp.ID.String(),
	}, nil
}
