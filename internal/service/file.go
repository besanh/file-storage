package service

import (
	"context"
	"time"

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

func (s *FileService) CreateFile(ctx context.Context, req *pb.CreateFileRequest) (*pb.CreateFileReply, error) {
	var parentIDPtr *uuid.UUID
	if req.ParentId != "" {
		parsedID, err := uuid.Parse(req.ParentId)
		if err != nil {
			return nil, err
		}
		parentIDPtr = &parsedID
	}

	resp, err := s.uc.CreateFile(ctx, biz.CreateFileRequest{
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

	return &pb.CreateFileReply{
		Id: resp.ID.String(),
	}, nil
}

func (s *FileService) GetUploadUrl(ctx context.Context, req *pb.GetUploadUrlRequest) (*pb.GetUploadUrlReply, error) {
	resp, err := s.uc.GetUploadUrl(ctx, biz.GetUploadUrlRequest{
		Name:         req.Name,
		FileSize:     req.FileSize,
		FileMimeType: req.FileMimeType,
		FileHash:     req.FileHash,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetUploadUrlReply{
		UploadUrl: resp.UploadUrl,
		FileId:    resp.FileID.String(),
	}, nil
}

func (s *FileService) GetDownloadUrl(ctx context.Context, req *pb.GetDownloadUrlRequest) (*pb.GetDownloadUrlReply, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	url, err := s.uc.GetDownloadUrl(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetDownloadUrlReply{
		DownloadUrl: url,
	}, nil
}

func (s *FileService) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileReply, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	resp, err := s.uc.GetFile(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetFileReply{
		Id:           resp.ID.String(),
		Name:         resp.Name,
		FileSize:     resp.FileSize,
		FileMimeType: resp.MimeType,
		IsFolder:     resp.IsFolder,
		CreatedAt:    resp.LastAccessed.Format(time.RFC3339),
		UpdatedAt:    resp.LastAccessed.Format(time.RFC3339),
	}, nil
}
