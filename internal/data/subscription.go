package data

import (
	"context"
	"database/sql"
	"file/internal/biz"
	db "file/internal/data/db/generated"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type subscriptionRepo struct {
	data *Data
	log  *log.Helper
}

func NewSubscriptionRepo(data *Data, logger log.Logger) biz.SubscriptionRepo {
	return &subscriptionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *subscriptionRepo) GetUserSubscription(ctx context.Context, userID uuid.UUID) (*biz.Subscription, error) {
	s, err := r.data.Queries(ctx).GetUserSubscription(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapDBSubscriptionToBiz(&s), nil
}

func (r *subscriptionRepo) GetSubscription(ctx context.Context, id uuid.UUID) (*biz.Subscription, error) {
	s, err := r.data.Queries(ctx).GetSubscription(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapDBSubscriptionToBiz(&s), nil
}

func (r *subscriptionRepo) ListUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]*biz.Subscription, error) {
	ss, err := r.data.Queries(ctx).ListUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]*biz.Subscription, 0, len(ss))
	for _, s := range ss {
		res = append(res, mapDBSubscriptionToBiz(&s))
	}
	return res, nil
}

func (r *subscriptionRepo) CreateSubscription(ctx context.Context, sub *biz.Subscription) (*biz.Subscription, error) {
	s, err := r.data.Queries(ctx).CreateSubscription(ctx, db.CreateSubscriptionParams{
		UserID:    sub.UserID,
		PlanID:    sub.PlanID,
		StartedAt: sub.StartedAt,
		ExpiredAt: sql.NullTime{Time: sub.ExpiredAt, Valid: !sub.ExpiredAt.IsZero()},
		Status:    sub.Status,
	})
	if err != nil {
		return nil, err
	}
	return mapDBSubscriptionToBiz(&s), nil
}

func (r *subscriptionRepo) UpdateSubscriptionStatus(ctx context.Context, id uuid.UUID, status string) (*biz.Subscription, error) {
	s, err := r.data.Queries(ctx).UpdateSubscriptionStatus(ctx, db.UpdateSubscriptionStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, err
	}
	return mapDBSubscriptionToBiz(&s), nil
}

func (r *subscriptionRepo) UpdateSubscriptionExpiration(ctx context.Context, id uuid.UUID, expiredAt time.Time) (*biz.Subscription, error) {
	s, err := r.data.Queries(ctx).UpdateSubscriptionExpiration(ctx, db.UpdateSubscriptionExpirationParams{
		ID:        id,
		ExpiredAt: sql.NullTime{Time: expiredAt, Valid: !expiredAt.IsZero()},
	})
	if err != nil {
		return nil, err
	}
	return mapDBSubscriptionToBiz(&s), nil
}

func (r *subscriptionRepo) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	return r.data.Queries(ctx).DeleteSubscription(ctx, id)
}

func mapDBSubscriptionToBiz(s *db.UserSubscription) *biz.Subscription {
	return &biz.Subscription{
		ID:        s.ID,
		UserID:    s.UserID,
		PlanID:    s.PlanID,
		StartedAt: s.StartedAt,
		ExpiredAt: s.ExpiredAt.Time,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt.Time,
	}
}
