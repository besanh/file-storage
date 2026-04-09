package biz

import (
	"context"
	"file/internal/common"
	"file/internal/conf"
	db "file/internal/data/db/generated"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type ShareRepo interface {
	CreateShareLink(ctx context.Context, s db.CreateShareLinkParams) (*db.ShareLink, error)
	GetShareLinkByToken(ctx context.Context, token string) (*db.ShareLink, error)
	RevokeShareLink(ctx context.Context, token string) error
	ListShareLinksByCreator(ctx context.Context, creator string, limit, offset int32) ([]*db.ShareLink, error)
	ListShareLinksByResource(ctx context.Context, resourceID uuid.UUID, limit, offset int32) ([]*db.ShareLink, error)
	UpdateShareLink(ctx context.Context, token string, permissionLevel string) error
}

type ShareUseCase struct {
	repo     ShareRepo
	log      *log.Helper
	authRepo AuthRepo
	tm       Transaction
	conf     *conf.Data
}

func NewShareUseCase(repo ShareRepo, logger log.Logger, authRepo AuthRepo, tm Transaction, conf *conf.Data) *ShareUseCase {
	return &ShareUseCase{
		repo:     repo,
		log:      log.NewHelper(logger),
		authRepo: authRepo,
		tm:       tm,
		conf:     conf,
	}
}

func (uc *ShareUseCase) CreateShareLink(ctx context.Context, input *ShareLinkRequest) (ShareLinkResponse, error) {
	actorType, actorID, err := GetActorInfo(ctx)
	if err != nil {
		return ShareLinkResponse{}, err
	}

	ownerID, err := uuid.Parse(actorID)
	if err != nil {
		if actorType == "service" {
			// For non-UUID client IDs, generate a deterministic UUID
			ownerID = uuid.NewMD5(uuid.NameSpaceOID, []byte(actorID))
		} else {
			return ShareLinkResponse{}, fmt.Errorf("invalid user identity string: %v", err)
		}
	}

	// PRE-FLIGHT CHECK: Is this person allowed to share this item?
	// We check the 'share' permission we defined in the SpiceDB schema!
	resp, err := uc.authRepo.CheckPermission(
		ctx,
		&CheckPermissionRequest{
			SubjectType: actorType,
			SubjectID:   actorID,
			Relation:    common.PermissionShare.String(),
			ObjectType:  common.TypeFile.String(),
			ObjectID:    input.ID.String(),
		},
	)
	if err != nil {
		return ShareLinkResponse{}, err
	}
	if !resp.Allowed {
		return ShareLinkResponse{}, fmt.Errorf("permission denied: cannot share file %s", input.ID.String())
	}

	// Generate a secure, URL-friendly token (21 characters)
	linkToken, err := gonanoid.New()
	if err != nil {
		return ShareLinkResponse{}, fmt.Errorf("failed to generate secure token: %w", err)
	}

	if err := uc.tm.ExecTx(ctx, func(ctx context.Context) error {
		if _, err := uc.repo.CreateShareLink(ctx, db.CreateShareLinkParams{
			LinkToken:       linkToken,
			ResourceID:      input.ID,
			ResourceType:    common.TypeFile.String(),
			CreatedBy:       ownerID.String(),
			PermissionLevel: input.PermissionLevel,
		}); err != nil {
			return err
		}

		// Grant Permission in SpiceDB
		if _, err := uc.authRepo.WriteRelationship(
			ctx,
			&WriteRelationshipRequest{
				ResourceType: common.TypeFile.String(),
				ResourceID:   input.ID.String(),
				Relation:     common.PermissionShare.String(),
				SubjectType:  actorType,
				SubjectID:    actorID,
			},
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		// Delete relationship if transaction failed
		_, _ = uc.authRepo.DeleteRelationship(ctx, &DeleteRelationshipRequest{
			ResourceType: common.TypeFile.String(),
			ResourceID:   input.ID.String(),
			Relation:     common.PermissionShare.String(),
			SubjectType:  actorType,
			SubjectID:    actorID,
		})
		return ShareLinkResponse{}, err
	}

	return ShareLinkResponse{
		Token: linkToken,
		URL:   uc.conf.Sharing.BaseUrl + linkToken,
	}, nil
}

func (uc *ShareUseCase) RevokeShareLink(ctx context.Context, input RevokeShareLinkRequest) (RevokeShareLinkResponse, error) {
	actorType, actorID, err := GetActorInfo(ctx)
	if err != nil {
		return RevokeShareLinkResponse{}, err
	}

	ownerID, err := uuid.Parse(actorID)
	if err != nil {
		if actorType == "service" {
			// For non-UUID client IDs, generate a deterministic UUID
			ownerID = uuid.NewMD5(uuid.NameSpaceOID, []byte(actorID))
		} else {
			return RevokeShareLinkResponse{}, fmt.Errorf("invalid user identity string: %v", err)
		}
	}

	// 1. Get the share link first to know which resource we are talking about
	shareLink, err := uc.repo.GetShareLinkByToken(ctx, input.Token)
	if err != nil {
		return RevokeShareLinkResponse{}, err
	}
	if shareLink == nil {
		return RevokeShareLinkResponse{}, fmt.Errorf("share link not found")
	}

	// PRE-FLIGHT CHECK: Does this person have permission to revoke (share) this item?
	resp, err := uc.authRepo.CheckPermission(
		ctx,
		&CheckPermissionRequest{
			SubjectType: actorType,
			SubjectID:   actorID,
			Relation:    common.PermissionShare.String(),
			ObjectType:  shareLink.ResourceType,
			ObjectID:    shareLink.ResourceID.String(),
		},
	)
	if err != nil {
		return RevokeShareLinkResponse{}, err
	}
	if !resp.Allowed {
		return RevokeShareLinkResponse{}, fmt.Errorf("permission denied: cannot revoke share link for file %s", shareLink.ResourceID.String())
	}

	if err := uc.tm.ExecTx(ctx, func(ctx context.Context) error {
		if shareLink.CreatedBy != ownerID.String() {
			return fmt.Errorf("permission denied: cannot revoke share link %s", input.Token)
		}
		if err := uc.repo.RevokeShareLink(ctx, input.Token); err != nil {
			return err
		}

		// Delete relationship in SpiceDB
		if _, err := uc.authRepo.DeleteRelationship(ctx, &DeleteRelationshipRequest{
			ResourceType: shareLink.ResourceType,
			ResourceID:   shareLink.ResourceID.String(),
			Relation:     common.PermissionShare.String(),
			SubjectType:  actorType,
			SubjectID:    actorID,
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		// Restore relationship if transaction failed
		_, _ = uc.authRepo.WriteRelationship(ctx, &WriteRelationshipRequest{
			ResourceType: shareLink.ResourceType,
			ResourceID:   shareLink.ResourceID.String(),
			Relation:     common.PermissionShare.String(),
			SubjectType:  actorType,
			SubjectID:    actorID,
		})
		return RevokeShareLinkResponse{}, err
	}

	return RevokeShareLinkResponse{}, nil
}

// func (uc *ShareUseCase) GetShareLink(ctx context.Context, token string) (*ShareLink, error) {
// 	return uc.repo.GetShareLinkByToken(ctx, token)
// }

func (uc *ShareUseCase) ListShareLinksByCreator(ctx context.Context, input ListShareLinksByCreatorRequest) (ListShareLinksByCreatorResponse, error) {
	actorType, actorID, err := GetActorInfo(ctx)
	if err != nil {
		return ListShareLinksByCreatorResponse{}, err
	}

	ownerID, err := uuid.Parse(actorID)
	if err != nil {
		if actorType == "service" {
			// For non-UUID client IDs, generate a deterministic UUID
			ownerID = uuid.NewMD5(uuid.NameSpaceOID, []byte(actorID))
		} else {
			return ListShareLinksByCreatorResponse{}, fmt.Errorf("invalid user identity string: %v", err)
		}
	}

	resp, err := uc.repo.ListShareLinksByCreator(ctx, ownerID.String(), input.Limit, input.Offset)
	if err != nil {
		return ListShareLinksByCreatorResponse{}, err
	}
	return ListShareLinksByCreatorResponse{
		ShareLinks: resp,
	}, nil
}

func (uc *ShareUseCase) ListShareLinksByResource(ctx context.Context, input ListShareLinksByResourceRequest) (ListShareLinksByResourceResponse, error) {
	actorType, actorID, err := GetActorInfo(ctx)
	if err != nil {
		return ListShareLinksByResourceResponse{}, err
	}

	// PRE-FLIGHT CHECK: Does this person have permission to view share links for this resource?
	resp, err := uc.authRepo.CheckPermission(
		ctx,
		&CheckPermissionRequest{
			SubjectType: actorType,
			SubjectID:   actorID,
			Relation:    common.PermissionShare.String(),
			ObjectType:  common.TypeFile.String(),
			ObjectID:    input.ID.String(),
		},
	)
	if err != nil {
		return ListShareLinksByResourceResponse{}, err
	}
	if !resp.Allowed {
		return ListShareLinksByResourceResponse{}, fmt.Errorf("permission denied: cannot list share links for resource %s", input.ID.String())
	}

	respList, err := uc.repo.ListShareLinksByResource(ctx, input.ID, input.Limit, input.Offset)
	if err != nil {
		return ListShareLinksByResourceResponse{}, err
	}
	return ListShareLinksByResourceResponse{
		ShareLinks: respList,
	}, nil
}

func (uc *ShareUseCase) UpdateUserPermission(ctx context.Context, input UpdateUserPermissionRequest) (UpdateUserPermissionResponse, error) {
	actorType, actorID, err := GetActorInfo(ctx)
	if err != nil {
		return UpdateUserPermissionResponse{}, err
	}

	ownerID, err := uuid.Parse(actorID)
	if err != nil {
		if actorType == "service" {
			// For non-UUID client IDs, generate a deterministic UUID
			ownerID = uuid.NewMD5(uuid.NameSpaceOID, []byte(actorID))
		} else {
			return UpdateUserPermissionResponse{}, fmt.Errorf("invalid user identity string: %v", err)
		}
	}

	// 1. Get link info to know the resource
	shareLink, err := uc.repo.GetShareLinkByToken(ctx, input.Token)
	if err != nil {
		return UpdateUserPermissionResponse{}, err
	}
	if shareLink == nil {
		return UpdateUserPermissionResponse{}, fmt.Errorf("share link not found")
	}

	// PRE-FLIGHT CHECK: Is this person allowed to share this item?
	chkResp, err := uc.authRepo.CheckPermission(
		ctx,
		&CheckPermissionRequest{
			SubjectType: actorType,
			SubjectID:   actorID,
			Relation:    common.PermissionShare.String(),
			ObjectType:  shareLink.ResourceType,
			ObjectID:    shareLink.ResourceID.String(),
		},
	)
	if err != nil {
		return UpdateUserPermissionResponse{}, err
	}
	if chkResp == nil || !chkResp.Allowed {
		return UpdateUserPermissionResponse{}, fmt.Errorf("permission denied: cannot update share link %s", input.Token)
	}

	if err := uc.tm.ExecTx(ctx, func(ctx context.Context) error {
		if shareLink.CreatedBy != ownerID.String() {
			return fmt.Errorf("permission denied: cannot update share link %s", input.Token)
		}
		if err := uc.repo.UpdateShareLink(ctx, input.Token, input.PermissionLevel); err != nil {
			return err
		}

		// Update relationship in SpiceDB
		if _, err := uc.authRepo.SwapRelationship(ctx, &SwapRelationshipRequest{
			ResourceType: shareLink.ResourceType,
			ResourceID:   shareLink.ResourceID.String(),
			SubjectType:  actorType,
			SubjectID:    actorID,
			OldRelation:  shareLink.PermissionLevel,
			NewRelation:  input.PermissionLevel,
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		// Restore relationship if transaction failed
		_, _ = uc.authRepo.WriteRelationship(ctx, &WriteRelationshipRequest{
			ResourceType: shareLink.ResourceType,
			ResourceID:   shareLink.ResourceID.String(),
			Relation:     common.PermissionShare.String(),
			SubjectType:  actorType,
			SubjectID:    actorID,
		})
		return UpdateUserPermissionResponse{}, err
	}

	return UpdateUserPermissionResponse{}, nil
}

func (uc *ShareUseCase) ResolveShareLink(ctx context.Context, input ResolveShareLinkRequest) (ResolveShareLinkResponse, error) {
	shareLink, err := uc.repo.GetShareLinkByToken(ctx, input.Token)
	if err != nil {
		return ResolveShareLinkResponse{}, err
	}
	if shareLink == nil {
		return ResolveShareLinkResponse{}, fmt.Errorf("share link not found")
	}

	// Check expiration
	if shareLink.ExpiresAt.Valid && shareLink.ExpiresAt.Time.Before(time.Now()) {
		return ResolveShareLinkResponse{}, fmt.Errorf("share link has expired")
	}

	return ResolveShareLinkResponse{
		ResourceID:      shareLink.ResourceID,
		ResourceType:    shareLink.ResourceType,
		PermissionLevel: shareLink.PermissionLevel,
		ExpiresAt:       shareLink.ExpiresAt.Time.Format(time.RFC3339),
		CreatorID:       shareLink.CreatedBy,
	}, nil
}
