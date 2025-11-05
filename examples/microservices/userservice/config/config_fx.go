package config

import "go.uber.org/fx"

var Module = fx.Module(
	"appconfigfx",

	fx.Provide(
		NewAppOptions,
	),
)
