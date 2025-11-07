package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/docs"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func (ic *UsersServiceConfigurator) configSwagger(routeBuilder contracts.RouteBuilder) {
	docs.SwaggerInfo.Version = "v1.0.0"
	docs.SwaggerInfo.Title = "User Service API"
	docs.SwaggerInfo.Description = "This is the User Service API documentation."

	routeBuilder.RegisterHandler(func(instance interface{}) {
		if app, ok := instance.(*fiber.App); ok {
			app.Get("/swagger/*", fiberSwagger.WrapHandler)
		}
	})
}
