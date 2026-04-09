package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type PlanUseCase struct {
	repo PlanRepo
	log  *log.Helper
}

func NewPlanUseCase(repo PlanRepo, logger log.Logger) *PlanUseCase {
	return &PlanUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *PlanUseCase) GetPlan(ctx context.Context, id uuid.UUID) (*Plan, error) {
	return uc.repo.GetPlan(ctx, id)
}

func (uc *PlanUseCase) GetPlanByName(ctx context.Context, name string) (*Plan, error) {
	return uc.repo.GetPlanByName(ctx, name)
}

func (uc *PlanUseCase) ListPlans(ctx context.Context) ([]*Plan, error) {
	return uc.repo.ListPlans(ctx)
}

func (uc *PlanUseCase) CreatePlan(ctx context.Context, plan *Plan) (*Plan, error) {
	return uc.repo.CreatePlan(ctx, plan)
}

func (uc *PlanUseCase) UpdatePlan(ctx context.Context, plan *Plan) (*Plan, error) {
	return uc.repo.UpdatePlan(ctx, plan)
}

func (uc *PlanUseCase) DeletePlan(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeletePlan(ctx, id)
}
