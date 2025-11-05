package fxparams

import (
	"github.com/go-playground/validator"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/logger"
	"go.uber.org/fx"
)

type UserRouteParams struct {
	fx.In

	Logger     logger.Logger
	UsersGroup contracts.RouteGroup `name:"user-routes" optional:"true"`
	Validator  *validator.Validate
}
