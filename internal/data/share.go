package data

import (
	"context"
	"database/sql"
	"file/internal/biz"
	db "file/internal/data/db/generated"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type shareRepo struct {
	data *Data
	log  *log.Helper
}

func NewShareRepo(data *Data, logger log.Logger) biz.ShareRepo {
	return &shareRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *shareRepo) CreateShareLink(ctx context.Context, s db.CreateShareLinkParams) (*db.ShareLink, error) {
	params := db.CreateShareLinkParams{
		LinkToken:       s.LinkToken,
		ResourceID:      s.ResourceID,
		ResourceType:    s.ResourceType,
		CreatedBy:       s.CreatedBy,
		PermissionLevel: s.PermissionLevel,
	}
	if !s.ExpiresAt.Valid {
		params.ExpiresAt = sql.NullTime{Time: s.ExpiresAt.Time, Valid: true}
	}

	resp, err := r.data.Queries(ctx).CreateShareLink(ctx, params)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *shareRepo) GetShareLinkByToken(ctx context.Context, token string) (*db.ShareLink, error) {
	resp, err := r.data.Queries(ctx).GetShareLinkByToken(ctx, token)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *shareRepo) RevokeShareLink(ctx context.Context, token string) error {
	return r.data.Queries(ctx).RevokeShareLink(ctx, token)
}

func (r *shareRepo) ListShareLinksByCreator(ctx context.Context, creator string, limit, offset int32) ([]*db.ShareLink, error) {
	models, err := r.data.Queries(ctx).ListShareLinksByCreator(ctx, db.ListShareLinksByCreatorParams{
		CreatedBy: creator,
		Limit:     limit,
		Offset:    offset,
	})
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	res := make([]*db.ShareLink, 0, len(models))
	for _, m := range models {
		res = append(res, &m)
	}
	return res, nil
}

func (r *shareRepo) ListShareLinksByResource(ctx context.Context, resourceID uuid.UUID, limit, offset int32) ([]*db.ShareLink, error) {
	models, err := r.data.Queries(ctx).ListShareLinksByResource(ctx, db.ListShareLinksByResourceParams{
		ResourceID: resourceID,
		Limit:      limit,
		Offset:     offset,
	})
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	res := make([]*db.ShareLink, 0, len(models))
	for _, m := range models {
		res = append(res, &m)
	}
	return res, nil
}

func (r *shareRepo) UpdateShareLink(ctx context.Context, token string, permissionLevel string) error {
	return r.data.Queries(ctx).UpdateShareLink(ctx, db.UpdateShareLinkParams{
		LinkToken:       token,
		PermissionLevel: permissionLevel,
	})
}
