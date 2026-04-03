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
	parentID, _ := uuid.Parse(req.ParentId)
	if err := s.uc.CreateFile(ctx, &parentID, req.Name, false, "", req.FileSize, req.FileType, req.FileExt, req.FileMimeType, req.FileVideoResolution, req.Status); err != nil {
		return nil, err
	}

	return &pb.CreateFileResponse{}, nil
}
