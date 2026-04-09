package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PlanID    uuid.UUID
	StartedAt time.Time
	ExpiredAt time.Time
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SubscriptionRepo interface {
	GetUserSubscription(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*Subscription, error)
	ListUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]*Subscription, error)
	CreateSubscription(ctx context.Context, sub *Subscription) (*Subscription, error)
	UpdateSubscriptionStatus(ctx context.Context, id uuid.UUID, status string) (*Subscription, error)
	UpdateSubscriptionExpiration(ctx context.Context, id uuid.UUID, expiredAt time.Time) (*Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
}
