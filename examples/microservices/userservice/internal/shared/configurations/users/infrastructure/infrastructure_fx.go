package infrastructure

import (
	customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
	"github.com/phatnt199/go-infra/pkg/core"
	postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
	"github.com/phatnt199/go-infra/pkg/migration/goose"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructurefx",

	core.Module,
	customfiber.Module,
	postgresgorm.Module,
	goose.Module,
)
