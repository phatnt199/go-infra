package config

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"authconfigfx",

	fx.Provide(
		NewAppOptions,
		NewAuthOptions,
		NewFiberHttpOptions,
		NewGormOptions,
		NewMigrationOptions,
		NewLogOptions,
	),
)
