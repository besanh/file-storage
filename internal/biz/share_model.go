package biz

import (
	db "file/internal/data/db/generated"

	"github.com/google/uuid"
)

type ShareLinkRequest struct {
	ID              uuid.UUID
	PermissionLevel string
}

type ShareLinkResponse struct {
	Token string
	URL   string
}

type RevokeShareLinkRequest struct {
	Token string
}

type RevokeShareLinkResponse struct {
	Success bool
}

type ListShareLinksByCreatorRequest struct {
	Limit  int32
	Offset int32
}

type ListShareLinksByCreatorResponse struct {
	ShareLinks []*db.ShareLink
}

type GetShareLinkByTokenRequest struct {
	Token string
}

type GetShareLinkByTokenResponse struct {
	ShareLink *db.ShareLink
}

type ListShareLinksByResourceRequest struct {
	ID     uuid.UUID
	Limit  int32
	Offset int32
}

type ListShareLinksByResourceResponse struct {
	ShareLinks []*db.ShareLink
}

type UpdateUserPermissionRequest struct {
	Token           string
	PermissionLevel string
}

type UpdateUserPermissionResponse struct {
}

type ResolveShareLinkRequest struct {
	Token string
}

type ResolveShareLinkResponse struct {
	ResourceID      uuid.UUID
	ResourceType    string
	PermissionLevel string
	ExpiresAt       string
	CreatorID       string
}
