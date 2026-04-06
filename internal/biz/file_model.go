package biz

import "github.com/google/uuid"

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
	ID *uuid.UUID
}
