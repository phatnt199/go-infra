package app

import (
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/configurations/users"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	"github.com/phatnt199/go-infra/pkg/logger"
	"go.uber.org/fx"
)

type UserApplication struct {
	*users.UsersServiceConfigurator
}

func NewUserApplication(
	provides []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *UserApplication {
	app := fxapp.NewApplication(provides, decorates, options, logger, environment)

	return &UserApplication{
		UsersServiceConfigurator: users.NewUsersServiceConfigurator(app),
	}
}
