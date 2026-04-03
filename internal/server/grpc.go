package server

import (
	"crypto/rsa"
	fileV1 "file/api/file/v1"
	v1 "file/api/helloworld/v1"
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
func NewGRPCServer(c *conf.Server, d *conf.Data, greeter *service.GreeterService, fileService *service.FileService, logger log.Logger, publicKey *rsa.PublicKey) *grpc.Server {
	// =========================================================================
	// 1. JWT Authentication Middleware
	// This uses the CPU to mathematically verify the token signature locally.
	// It requires ZERO network calls to the Auth Service.
	// =========================================================================
	jwtAuthn := jwt.Server(func(token *jwtv5.Token) (any, error) {
		return publicKey, nil
	}, jwt.WithSigningMethod(jwtv5.SigningMethodRS256))

	// =========================================================================
	// 2. Route Selector
	// This decides which gRPC endpoints require the JWT token and which are open.
	// =========================================================================
	authSelector := selector.Server(jwtAuthn).Match(NewWhiteListMatcher([]string{
		"/api.file.v1.FileService/HealthCheck", // Skip Health Checks
		"/health",                              // General health check
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
	v1.RegisterGreeterServer(srv, greeter)
	fileV1.RegisterFileServiceServer(srv, fileService)
	return srv
}
