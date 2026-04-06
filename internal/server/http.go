package server

import (
	"crypto/rsa"
	fileV1 "file/api/file/v1"
	"file/internal/conf"
	"file/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, d *conf.Data, fileService *service.FileService, logger log.Logger, publicKey *rsa.PublicKey) *http.Server {
	jwtAuthn := jwt.Server(func(token *jwtv5.Token) (any, error) {
		return publicKey, nil
	}, jwt.WithSigningMethod(jwtv5.SigningMethodRS256))

	authSelector := selector.Server(jwtAuthn).Match(NewWhiteListMatcher([]string{
		"/api.file.v1.FileService/HealthCheck",
		"/health",
	})).Build()

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			authSelector,
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	fileV1.RegisterFileServiceHTTPServer(srv, fileService)
	return srv
}
