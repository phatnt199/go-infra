package data

import (
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/data/dbcontext"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"datafx",

	fx.Provide(
		dbcontext.NewAuthGormDBContext,
	),
)
