package config

import (
	"github.com/phatnt199/go-infra/pkg/crypto"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"appconfigfx",

	fx.Provide(
		NewAppOptions,
	),

	fx.Provide(
		func(opts *JWTOptions) *crypto.JWTConfig {
			return opts.ToJWTConfig()
		},
		NewJWTOptions,
	),
)
