package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type SubscriptionUseCase struct {
	repo     SubscriptionRepo
	planRepo PlanRepo
	log      *log.Helper
}

func NewSubscriptionUseCase(repo SubscriptionRepo, planRepo PlanRepo, logger log.Logger) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		repo:     repo,
		planRepo: planRepo,
		log:      log.NewHelper(logger),
	}
}

func (uc *SubscriptionUseCase) GetUserSubscription(ctx context.Context, userID uuid.UUID) (*Subscription, error) {
	return uc.repo.GetUserSubscription(ctx, userID)
}

func (uc *SubscriptionUseCase) SubscribePlan(ctx context.Context, userID uuid.UUID, planID uuid.UUID) (*Subscription, error) {
	plan, err := uc.planRepo.GetPlan(ctx, planID)
	if err != nil {
		return nil, err
	}

	startedAt := time.Now()
	expiredAt := startedAt.AddDate(0, 0, int(plan.DurationDays))

	sub := &Subscription{
		UserID:    userID,
		PlanID:    planID,
		StartedAt: startedAt,
		ExpiredAt: expiredAt,
		Status:    "active",
	}

	return uc.repo.CreateSubscription(ctx, sub)
}

func (uc *SubscriptionUseCase) ListUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]*Subscription, error) {
	return uc.repo.ListUserSubscriptions(ctx, userID)
}
