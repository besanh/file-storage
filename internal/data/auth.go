package data

import (
	"context"
	m2mv1 "file/api/m2m_auth/v1"
	"file/internal/biz"
	"file/internal/conf"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type authRepo struct {
	m2mClient    m2mv1.AuthClient
	clientID     string
	clientSecret string
	kid          string
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

	return r, nil
}

func (r *authRepo) CheckPermission(ctx context.Context, subjectType, subjectID, relation, objectType, objectID string) (bool, error) {
	// For now, return true to allow debugging
	return true, nil
}
