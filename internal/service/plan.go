package service

import (
	"context"
	pb "file/api/plan/v1"
	"file/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type PlanService struct {
	pb.UnimplementedPlanServiceServer
	uc  *biz.PlanUseCase
	log *log.Helper
}

func NewPlanService(uc *biz.PlanUseCase, logger log.Logger) *PlanService {
	return &PlanService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *PlanService) GetPlan(ctx context.Context, req *pb.GetPlanRequest) (*pb.GetPlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	p, err := s.uc.GetPlan(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.GetPlanResponse{Plan: mapBizPlanToPb(p)}, nil
}

func (s *PlanService) ListPlans(ctx context.Context, req *pb.ListPlansRequest) (*pb.ListPlansResponse, error) {
	ps, err := s.uc.ListPlans(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Plan, 0, len(ps))
	for _, p := range ps {
		res = append(res, mapBizPlanToPb(p))
	}
	return &pb.ListPlansResponse{Plans: res}, nil
}

func (s *PlanService) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest) (*pb.CreatePlanResponse, error) {
	p, err := s.uc.CreatePlan(ctx, &biz.Plan{
		Name:          req.Name,
		StorageQuota:  req.StorageQuota,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		DurationDays:  req.DurationDays,
		Description:   req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreatePlanResponse{Plan: mapBizPlanToPb(p)}, nil
}

func (s *PlanService) UpdatePlan(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.UpdatePlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	p, err := s.uc.UpdatePlan(ctx, &biz.Plan{
		ID:            id,
		Name:          req.Name,
		StorageQuota:  req.StorageQuota,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		DurationDays:  req.DurationDays,
		Description:   req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePlanResponse{Plan: mapBizPlanToPb(p)}, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, req *pb.DeletePlanRequest) (*pb.DeletePlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	err = s.uc.DeletePlan(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.DeletePlanResponse{Success: true}, nil
}

func mapBizPlanToPb(p *biz.Plan) *pb.Plan {
	if p == nil {
		return nil
	}
	return &pb.Plan{
		Id:            p.ID.String(),
		Name:          p.Name,
		StorageQuota:  p.StorageQuota,
		Price:         p.Price,
		DiscountPrice: p.DiscountPrice,
		DurationDays:  p.DurationDays,
		Description:   p.Description,
		CreatedAt:     p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     p.UpdatedAt.Format(time.RFC3339),
	}
}
