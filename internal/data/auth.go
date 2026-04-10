package data

import (
	"context"
	m2mv1 "file/api/m2m_auth/v1"
	permissionV1 "file/api/permission/v1"
	userV1 "file/api/user/v1"
	"file/internal/biz"
	"file/internal/conf"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type authRepo struct {
	m2mClient        m2mv1.M2MAuthServiceClient
	clientID         string
	clientSecret     string
	kid              string
	permissionClient permissionV1.PermissionServiceClient
	userClient       userV1.UserServiceClient
}

func NewAuthRepo(c *conf.Data) (biz.AuthRepo, error) {
	if c.Auth == nil {
		return nil, fmt.Errorf("auth configuration is missing")
	}

	r := &authRepo{
		clientID:     c.Auth.ClientId,
		clientSecret: c.Auth.ClientSecret,
		kid:          c.Auth.Kid,
	}

	// Setup gRPC connection with middleware
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.Auth.Host),
		grpc.WithMiddleware(
			recovery.Recovery(),
			circuitbreaker.Client(),
		),
	)
	if err != nil {
		return nil, err
	}

	r.m2mClient = m2mv1.NewM2MAuthServiceClient(conn)
	r.permissionClient = permissionV1.NewPermissionServiceClient(conn)
	r.userClient = userV1.NewUserServiceClient(conn)

	return r, nil
}

func (r *authRepo) CheckPermission(ctx context.Context, req *biz.CheckPermissionRequest) (*biz.CheckPermissionResponse, error) {
	resp, err := r.permissionClient.CheckPermission(ctx, &permissionV1.CheckPermissionRequest{
		SubjectType:  req.SubjectType,
		SubjectId:    req.SubjectID,
		Relation:     req.Relation,
		ResourceType: req.ObjectType,
		ResourceId:   req.ObjectID,
	})
	if err != nil {
		return nil, err
	}
	return &biz.CheckPermissionResponse{Allowed: resp.Allowed}, nil
}

func (r *authRepo) WriteRelationship(ctx context.Context, req *biz.WriteRelationshipRequest) (*biz.WriteRelationshipResponse, error) {
	_, err := r.permissionClient.WriteRelationship(ctx, &permissionV1.WriteRelationshipRequest{
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceID,
		Relation:     req.Relation,
		SubjectType:  req.SubjectType,
		SubjectId:    req.SubjectID,
	})
	if err != nil {
		return nil, err
	}
	return &biz.WriteRelationshipResponse{}, nil
}

func (r *authRepo) DeleteRelationship(ctx context.Context, req *biz.DeleteRelationshipRequest) (*biz.DeleteRelationshipResponse, error) {
	_, err := r.permissionClient.DeleteRelationship(ctx, &permissionV1.DeleteRelationshipRequest{
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceID,
		Relation:     req.Relation,
		SubjectType:  req.SubjectType,
		SubjectId:    req.SubjectID,
	})
	if err != nil {
		return nil, err
	}
	return &biz.DeleteRelationshipResponse{}, nil
}

func (r *authRepo) SwapRelationship(ctx context.Context, req *biz.SwapRelationshipRequest) (*biz.SwapRelationshipResponse, error) {
	_, err := r.permissionClient.SwapRelationship(ctx, &permissionV1.SwapRelationshipRequest{
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceID,
		SubjectType:  req.SubjectType,
		SubjectId:    req.SubjectID,
		OldRelation:  req.OldRelation,
		NewRelation:  req.NewRelation,
	})
	if err != nil {
		return nil, err
	}
	return &biz.SwapRelationshipResponse{}, nil
}

func (r *authRepo) GetUserProfile(ctx context.Context, id string) (*biz.UserProfile, error) {
	resp, err := r.userClient.GetUser(ctx, &userV1.GetUserRequest{UserId: id})
	if err != nil {
		return nil, err
	}
	return &biz.UserProfile{
		UserID:    resp.UserId,
		Email:     resp.Email,
		Role:      resp.Role,
		Scope:     resp.Scope,
		Status:    resp.Status,
		CreatedAt: resp.CreatedAt,
	}, nil
}
