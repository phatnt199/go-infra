package configurations

import (
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/configurations/endpoints"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp/contracts"
	"github.com/phatnt199/go-infra/pkg/core/web/route"
	"go.uber.org/fx"
)

type UsersModuleConfigurator struct {
	contracts.Application
}

type MapUsersEndpointParams struct {
	fx.In
	AuthEndpoints []route.Endpoint `group:"auth-routes"`
	UserEndpoints []route.Endpoint `group:"user-routes"`
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
	umc.ResolveFunc(
		func(params MapUsersEndpointParams) error {
			allEndpoints := append(params.AuthEndpoints, params.UserEndpoints...)
			return endpoints.RegisterEndpoints(allEndpoints)
		},
	)
	return nil
}
