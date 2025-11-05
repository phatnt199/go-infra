package users

import (
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/repositories"
	createuserv1 "github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/features/createuser/v1"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/core/web/route"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"usersfx",

	fx.Provide(repositories.NewPostgresUserRepository),

	fx.Provide(
		fx.Annotate(func(usersServer contracts.HttpServer) contracts.RouteGroup {
			v1 := usersServer.RouteBuilder().Group("/api/v1")
			return v1.Group("/users")
		}, fx.ResultTags(`name:"user-http-group"`))),

	fx.Provide(
		route.AsRoute(
			createuserv1.NewCreateUserHandler,
			"user-routes",
		),
	),
)
