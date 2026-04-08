package data

import (
	"context"
	m2mv1 "file/api/m2m_auth/v1"
	permissionV1 "file/api/permission/v1"
	"file/internal/biz"
	"file/internal/conf"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type authRepo struct {
	m2mClient        m2mv1.AuthClient
	clientID         string
	clientSecret     string
	kid              string
	permissionClient permissionV1.PermissionClient
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

	r.m2mClient = m2mv1.NewAuthClient(conn)
	r.permissionClient = permissionV1.NewPermissionClient(conn)

	return r, nil
}

func (r *authRepo) CheckPermission(ctx context.Context, subjectType, subjectID, relation, objectType, objectID string) (bool, error) {
	resp, err := r.permissionClient.CheckPermission(ctx, &permissionV1.CheckPermissionRequest{
		SubjectType:  subjectType,
		SubjectId:    subjectID,
		Relation:     relation,
		ResourceType: objectType,
		ResourceId:   objectID,
	})
	if err != nil {
		return false, err
	}
	return resp.Allowed, nil
}

func (r *authRepo) WriteRelationship(ctx context.Context, resourceType, resourceID, relation, subjectType, subjectID string) error {
	_, err := r.permissionClient.WriteRelationship(ctx, &permissionV1.WriteRelationshipRequest{
		ResourceType: resourceType,
		ResourceId:   resourceID,
		Relation:     relation,
		SubjectType:  subjectType,
		SubjectId:    subjectID,
	})
	return err
}

func (r *authRepo) DeleteRelationship(ctx context.Context, resourceType, resourceID, relation, subjectType, subjectID string) error {
	_, err := r.permissionClient.DeleteRelationship(ctx, &permissionV1.DeleteRelationshipRequest{
		ResourceType: resourceType,
		ResourceId:   resourceID,
		Relation:     relation,
		SubjectType:  subjectType,
		SubjectId:    subjectID,
	})
	return err
}
