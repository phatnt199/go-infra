package gomigrate

import (
	"github.com/phatnt199/go-infra/pkg/migration"

	"go.uber.org/fx"
)

var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module(
		"gomigratefx",
		mongoProviders,
	)

	mongoProviders = fx.Provide(
		migration.ProvideConfig,
		NewGoMigratorPostgres,
	)
)
