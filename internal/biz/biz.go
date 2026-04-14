package biz

import (
	"errors"

	"github.com/google/wire"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewFileUsecase, NewPublicKey, NewPhysicalFileUsecase, NewShareUseCase, NewPlanUseCase, NewSubscriptionUseCase, NewDashboardUsecase, NewEventBus)
