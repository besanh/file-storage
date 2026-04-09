package service

import (
	"context"
	pb "file/api/share/v1"
	"file/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type ShareService struct {
	pb.UnimplementedShareServiceServer
	suc *biz.ShareUseCase
	log *log.Helper
}

func NewShareService(suc *biz.ShareUseCase, logger log.Logger) *ShareService {
	return &ShareService{
		suc: suc,
		log: log.NewHelper(logger),
	}
}

func (s *ShareService) CreateShareLink(ctx context.Context, req *pb.CreateShareLinkRequest) (*pb.CreateShareLinkResponse, error) {
	resourceID, err := uuid.Parse(req.ResourceId)
	if err != nil {
		return nil, err
	}

	res, err := s.suc.CreateShareLink(ctx, &biz.ShareLinkRequest{
		ID:              resourceID,
		PermissionLevel: req.PermissionLevel,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateShareLinkResponse{
		LinkToken: res.Token,
		ShareUrl:  res.URL,
	}, nil
}

func (s *ShareService) ResolveShareLink(ctx context.Context, req *pb.ResolveShareLinkRequest) (*pb.ResolveShareLinkResponse, error) {
	res, err := s.suc.ResolveShareLink(ctx, biz.ResolveShareLinkRequest{
		Token: req.LinkToken,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ResolveShareLinkResponse{
		ResourceId:      res.ResourceID.String(),
		ResourceType:    res.ResourceType,
		PermissionLevel: res.PermissionLevel,
		ExpiresAt:       res.ExpiresAt,
		CreatorId:       res.CreatorID,
	}, nil
}

func (s *ShareService) RevokeShareLink(ctx context.Context, req *pb.RevokeShareLinkRequest) (*pb.RevokeShareLinkResponse, error) {
	resp, err := s.suc.RevokeShareLink(ctx, biz.RevokeShareLinkRequest{
		Token: req.LinkToken,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RevokeShareLinkResponse{
		Success: resp.Success,
	}, nil
}

func (s *ShareService) ListShareLinksByCreator(ctx context.Context, req *pb.ListShareLinksByCreatorRequest) (*pb.ListShareLinksByCreatorResponse, error) {
	res, err := s.suc.ListShareLinksByCreator(ctx, biz.ListShareLinksByCreatorRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}

	links := make([]*pb.ShareLinkInfo, 0, len(res.ShareLinks))
	for _, link := range res.ShareLinks {
		links = append(links, &pb.ShareLinkInfo{
			LinkToken:       link.LinkToken,
			PermissionLevel: link.PermissionLevel,
			CreatedAt:       link.CreatedAt.Time.Format(time.RFC3339),
			ShareUrl:        "https://besanh.com/s/" + link.LinkToken,
		})
	}

	return &pb.ListShareLinksByCreatorResponse{
		Links: links,
	}, nil
}

func (s *ShareService) ListShareLinksByResource(ctx context.Context, req *pb.ListShareLinksByResourceRequest) (*pb.ListShareLinksByResourceResponse, error) {
	resourceID, err := uuid.Parse(req.ResourceId)
	if err != nil {
		return nil, err
	}

	res, err := s.suc.ListShareLinksByResource(ctx, biz.ListShareLinksByResourceRequest{
		ID:     resourceID,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}

	links := make([]*pb.ShareLinkInfo, 0, len(res.ShareLinks))
	for _, link := range res.ShareLinks {
		links = append(links, &pb.ShareLinkInfo{
			LinkToken:       link.LinkToken,
			PermissionLevel: link.PermissionLevel,
			CreatedAt:       link.CreatedAt.Time.Format(time.RFC3339),
			ShareUrl:        "https://besanh.com/s/" + link.LinkToken,
		})
	}

	return &pb.ListShareLinksByResourceResponse{
		Links: links,
	}, nil
}

func (s *ShareService) UpdateShareLink(ctx context.Context, req *pb.UpdateShareLinkRequest) (*pb.UpdateShareLinkResponse, error) {
	_, err := s.suc.UpdateUserPermission(ctx, biz.UpdateUserPermissionRequest{
		Token:           req.LinkToken,
		PermissionLevel: req.NewPermission,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateShareLinkResponse{
		Success: true,
	}, nil
}
