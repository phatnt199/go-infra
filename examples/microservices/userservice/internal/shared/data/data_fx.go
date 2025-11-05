package data

import (
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/data/dbcontext"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"datafx",

	fx.Provide(
		dbcontext.NewUsersGormDBContext,
	),
)
