package users

import (
	"fmt"
	"net/http"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/config"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/configurations/users/infrastructure"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/configurations"
	fxcontracts "github.com/phatnt199/go-infra/pkg/adapter/fxapp/contracts"
	httpcontracts "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	migrationcontracts "github.com/phatnt199/go-infra/pkg/migration/contracts"
	"gorm.io/gorm"
)

type UsersServiceConfigurator struct {
	fxcontracts.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	usersModuleConfigurator    *configurations.UsersModuleConfigurator
}

func NewUsersServiceConfigurator(
	fxapp fxcontracts.Application,
) *UsersServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(fxapp)
	usersModuleConfigurator := configurations.NewUsersModuleConfigurator(fxapp)

	return &UsersServiceConfigurator{
		Application:                fxapp,
		infrastructureConfigurator: infraConfigurator,
		usersModuleConfigurator:    usersModuleConfigurator,
	}
}

func (usc *UsersServiceConfigurator) ConfigureUsersService() error {
	// Shared | Infrastructure

	usc.infrastructureConfigurator.ConfigureInfrastructure()

	// Shared | Users configurations
	usc.ResolveFunc(
		func(db *gorm.DB, postgresMigrationRunner migrationcontracts.PostgresMigrationRunner) error {
			err := usc.migrateUsers(postgresMigrationRunner)
			if err != nil {
				return err
			}
			err = usc.seedUsers(db)
			if err != nil {
				return err
			}
			return nil
		},
	)

	err := usc.usersModuleConfigurator.ConfigureUsersModule()
	return err
}

func (usc *UsersServiceConfigurator) MapUsersEndpoints() error {
	usc.ResolveFunc(
		func(usersServer httpcontracts.HttpServer, options *config.AppOptions) error {
			usersServer.SetupDefaultMiddlewares()

			usersServer.RouteBuilder().GET("", func(ctx httpcontracts.Context) error {
				return ctx.JSON(http.StatusOK, fmt.Sprintf("%s is running!", options.GetMicroserviceNameUpper()))
			})

			usc.configSwagger(usersServer.RouteBuilder())

			return nil
		})

	err := usc.usersModuleConfigurator.MapUsersEndpoint()

	return err
}
