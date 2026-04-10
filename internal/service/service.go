package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService, NewFileService, NewShareService, NewPlanService, NewSubscriptionService, NewDashboardService)
