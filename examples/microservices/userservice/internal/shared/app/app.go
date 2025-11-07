package app

import "github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/configurations/users"

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	appBuilder := NewUserApplicationBuilder()
	appBuilder.ProvideModule(users.UsersServiceModule)

	app := appBuilder.Build()

	err := app.ConfigureUsersService()
	if err != nil {
		app.Logger().Fatal("Failed to configure users service", "error", err)
	}

	err = app.MapUsersEndpoints()
	if err != nil {
		app.Logger().Fatal("Failed to map users endpoints", "error", err)
	}

	app.Logger().Info("User service is running...")

	app.Run()
}
