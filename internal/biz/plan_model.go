package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID            uuid.UUID
	Name          string
	StorageQuota  int64
	Price         int64
	DiscountPrice int64
	DurationDays  int32
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PlanRepo interface {
	GetPlan(ctx context.Context, id uuid.UUID) (*Plan, error)
	GetPlanByName(ctx context.Context, name string) (*Plan, error)
	ListPlans(ctx context.Context) ([]*Plan, error)
	CreatePlan(ctx context.Context, plan *Plan) (*Plan, error)
	UpdatePlan(ctx context.Context, plan *Plan) (*Plan, error)
	DeletePlan(ctx context.Context, id uuid.UUID) error
}
