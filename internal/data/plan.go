package data

import (
	"context"
	"database/sql"
	"file/internal/biz"
	db "file/internal/data/db/generated"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type planRepo struct {
	data *Data
	log  *log.Helper
}

func NewPlanRepo(data *Data, logger log.Logger) biz.PlanRepo {
	return &planRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *planRepo) GetPlan(ctx context.Context, id uuid.UUID) (*biz.Plan, error) {
	p, err := r.data.Queries(ctx).GetPlan(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapDBPlanToBiz(&p), nil
}

func (r *planRepo) GetPlanByName(ctx context.Context, name string) (*biz.Plan, error) {
	p, err := r.data.Queries(ctx).GetPlanByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapDBPlanToBiz(&p), nil
}

func (r *planRepo) ListPlans(ctx context.Context) ([]*biz.Plan, error) {
	ps, err := r.data.Queries(ctx).ListPlans(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*biz.Plan, 0, len(ps))
	for _, p := range ps {
		res = append(res, mapDBPlanToBiz(&p))
	}
	return res, nil
}

func (r *planRepo) CreatePlan(ctx context.Context, plan *biz.Plan) (*biz.Plan, error) {
	p, err := r.data.Queries(ctx).CreatePlan(ctx, db.CreatePlanParams{
		Name:          plan.Name,
		StorageQuota:  plan.StorageQuota,
		Price:         plan.Price,
		DiscountPrice: plan.DiscountPrice,
		DurationDays:  plan.DurationDays,
		Description:   sql.NullString{String: plan.Description, Valid: plan.Description != ""},
	})
	if err != nil {
		return nil, err
	}
	return mapDBPlanToBiz(&p), nil
}

func (r *planRepo) UpdatePlan(ctx context.Context, plan *biz.Plan) (*biz.Plan, error) {
	p, err := r.data.Queries(ctx).UpdatePlan(ctx, db.UpdatePlanParams{
		ID:            plan.ID,
		Name:          plan.Name,
		StorageQuota:  plan.StorageQuota,
		Price:         plan.Price,
		DiscountPrice: plan.DiscountPrice,
		DurationDays:  plan.DurationDays,
		Description:   sql.NullString{String: plan.Description, Valid: plan.Description != ""},
	})
	if err != nil {
		return nil, err
	}
	return mapDBPlanToBiz(&p), nil
}

func (r *planRepo) DeletePlan(ctx context.Context, id uuid.UUID) error {
	return r.data.Queries(ctx).DeletePlan(ctx, id)
}

func mapDBPlanToBiz(p *db.Plan) *biz.Plan {
	return &biz.Plan{
		ID:            p.ID,
		Name:          p.Name,
		StorageQuota:  p.StorageQuota,
		Price:         p.Price,
		DiscountPrice: p.DiscountPrice,
		DurationDays:  p.DurationDays,
		Description:   p.Description.String,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt.Time,
	}
}
