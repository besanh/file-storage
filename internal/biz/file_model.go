package biz

import (
	"time"

	"github.com/google/uuid"
)

type CreateFolderRequest struct {
	ParentID *uuid.UUID
	Name     string
}

type CreateFolderResponse struct {
	ID *uuid.UUID
}

type CreateFileRequest struct {
	ParentID            *uuid.UUID
	Name                string
	IsFolder            bool
	FileHash            string
	FileSize            int64
	FileType            string
	FileExt             string
	FileMimeType        string
	FileVideoResolution string
	Status              string
}

type CreateFileResponse struct {
	ID uuid.UUID
}

type FileNode struct {
	ID           uuid.UUID
	Name         string
	IsFolder     bool
	FileSize     int64
	FileType     string
	FileExt      string
	MimeType     string
	LastAccessed time.Time
}

type GetUploadUrlRequest struct {
	Name         string
	FileSize     int64
	FileMimeType string
	FileHash     string
}

type GetUploadUrlResponse struct {
	UploadUrl string
	FileID    uuid.UUID
}
