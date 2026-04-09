package service

import (
	"context"
	pb "file/api/subscription/v1"
	"file/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	pb.UnimplementedSubscriptionServiceServer
	uc  *biz.SubscriptionUseCase
	log *log.Helper
}

func NewSubscriptionService(uc *biz.SubscriptionUseCase, logger log.Logger) *SubscriptionService {
	return &SubscriptionService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *SubscriptionService) GetUserSubscription(ctx context.Context, req *pb.GetUserSubscriptionRequest) (*pb.GetUserSubscriptionResponse, error) {
	userIDStr := biz.GetUserID(ctx)
	if userIDStr == "" {
		return nil, biz.ErrUnauthorized
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	sub, err := s.uc.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserSubscriptionResponse{Subscription: mapBizSubscriptionToPb(sub)}, nil
}

func (s *SubscriptionService) SubscribePlan(ctx context.Context, req *pb.SubscribePlanRequest) (*pb.SubscribePlanResponse, error) {
	userIDStr := biz.GetUserID(ctx)
	if userIDStr == "" {
		return nil, biz.ErrUnauthorized
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	planID, err := uuid.Parse(req.PlanId)
	if err != nil {
		return nil, err
	}

	sub, err := s.uc.SubscribePlan(ctx, userID, planID)
	if err != nil {
		return nil, err
	}
	return &pb.SubscribePlanResponse{Subscription: mapBizSubscriptionToPb(sub)}, nil
}

func (s *SubscriptionService) ListUserSubscriptions(ctx context.Context, req *pb.ListUserSubscriptionsRequest) (*pb.ListUserSubscriptionsResponse, error) {
	userIDStr := biz.GetUserID(ctx)
	if userIDStr == "" {
		return nil, biz.ErrUnauthorized
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	ss, err := s.uc.ListUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Subscription, 0, len(ss))
	for _, sub := range ss {
		res = append(res, mapBizSubscriptionToPb(sub))
	}
	return &pb.ListUserSubscriptionsResponse{Subscriptions: res}, nil
}

func mapBizSubscriptionToPb(s *biz.Subscription) *pb.Subscription {
	if s == nil {
		return nil
	}
	return &pb.Subscription{
		Id:        s.ID.String(),
		UserId:    s.UserID.String(),
		PlanId:    s.PlanID.String(),
		StartedAt: s.StartedAt.Format(time.RFC3339),
		ExpiredAt: s.ExpiredAt.Format(time.RFC3339),
		Status:    s.Status,
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
	}
}
