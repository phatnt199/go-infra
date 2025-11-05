package configurations

import (
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/configurations/endpoints"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp/contracts"
)

type UsersModuleConfigurator struct {
	contracts.Application
}

func NewUsersModuleConfigurator(
	fxapp contracts.Application,
) *UsersModuleConfigurator {
	return &UsersModuleConfigurator{
		Application: fxapp,
	}
}

func (umc *UsersModuleConfigurator) ConfigureUsersModule() error {
	return nil
}

func (umc *UsersModuleConfigurator) MapUsersEndpoint() error {
	umc.ResolveFuncWithParamTag(
		endpoints.RegisterEndpoints,
		`group:"user-handlers"`,
	)
	return nil
}
