package server

import (
	"crypto/rsa"
	fileV1 "file/api/file/v1"
	planV1 "file/api/plan/v1"
	shareV1 "file/api/share/v1"
	subV1 "file/api/subscription/v1"
	"file/internal/conf"
	"file/internal/service"

	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, d *conf.Data, fileService *service.FileService, shareService *service.ShareService, planService *service.PlanService, subService *service.SubscriptionService, logger log.Logger, publicKey *rsa.PublicKey) *grpc.Server {
	jwtAuthn := jwt.Server(func(token *jwtv5.Token) (any, error) {
		return publicKey, nil
	}, jwt.WithSigningMethod(jwtv5.SigningMethodRS256))

	authSelector := selector.Server(jwtAuthn).Match(NewWhiteListMatcher([]string{
		"/api.file.v1.FileService/HealthCheck",
		"/api.plan.v1.PlanService/ListPlans",
		"/api.plan.v1.PlanService/GetPlan",
		"/health",
	})).Build()

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			authSelector,
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	fileV1.RegisterFileServiceServer(srv, fileService)
	shareV1.RegisterShareServiceServer(srv, shareService)
	planV1.RegisterPlanServiceServer(srv, planService)
	subV1.RegisterSubscriptionServiceServer(srv, subService)

	return srv
}
