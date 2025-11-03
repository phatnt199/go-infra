package postgresgorm

import (
	"github.com/phatnt199/go-infra/pkg/health/contracts"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"gormpostgresfx",

	fx.Provide(
		provideConfig,
		NewGorm,

		fx.Annotate(
			NewGormHealthChecker,
			fx.As(new(contracts.Health)),
			fx.ResultTags(`group:"healths"`),
		),
	),
)
