package server

import (
	"crypto/rsa"
	dashboardV1 "file/api/dashboard/v1"
	fileV1 "file/api/file/v1"
	greeterV1 "file/api/helloworld/v1"
	planV1 "file/api/plan/v1"
	shareV1 "file/api/share/v1"
	subV1 "file/api/subscription/v1"
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
func NewHTTPServer(c *conf.Server, d *conf.Data, greeterService *service.GreeterService, fileService *service.FileService, shareService *service.ShareService, planService *service.PlanService, subService *service.SubscriptionService, dashService *service.DashboardService, logger log.Logger, publicKey *rsa.PublicKey) *http.Server {
	jwtAuthn := jwt.Server(func(token *jwtv5.Token) (any, error) {
		return publicKey, nil
	}, jwt.WithSigningMethod(jwtv5.SigningMethodRS256))

	authSelector := selector.Server(jwtAuthn).Match(NewWhiteListMatcher([]string{
		"/api.file.v1.FileService/HealthCheck",
		"/api.plan.v1.PlanService/ListPlans",
		"/api.plan.v1.PlanService/GetPlan",
		"/health",
		"/api.helloworld.v1.Greeter/SayHello",
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
	greeterV1.RegisterGreeterHTTPServer(srv, greeterService)
	fileV1.RegisterFileServiceHTTPServer(srv, fileService)
	shareV1.RegisterShareServiceHTTPServer(srv, shareService)
	planV1.RegisterPlanServiceHTTPServer(srv, planService)
	subV1.RegisterSubscriptionServiceHTTPServer(srv, subService)
	dashboardV1.RegisterDashboardServiceHTTPServer(srv, dashService)

	// Custom route for SSE
	srv.HandleFunc("/v1/dashboard/updates", dashService.StreamDashboardUpdates)

	return srv
}
