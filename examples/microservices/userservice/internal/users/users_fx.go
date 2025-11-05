package users

import (
	"github.com/go-playground/validator"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/repositories"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/features/auth"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/features/profile"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/services"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/core/web/route"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"usersfx",

	// Provide validator
	fx.Provide(func() *validator.Validate {
		return validator.New()
	}),

	// Provide repository
	fx.Provide(repositories.NewPostgresUserRepository),

	// Provide services
	fx.Provide(services.NewUserService),

	// Provide route groups
	fx.Provide(
		fx.Annotate(func(usersServer contracts.HttpServer) contracts.RouteGroup {
			basePath := usersServer.Cfg().GetBasePath()
			v1 := usersServer.RouteBuilder().Group(basePath)
			return v1.Group("/users")
		}, fx.ResultTags(`name:"user-routes"`))),

	fx.Provide(
		fx.Annotate(func(usersServer contracts.HttpServer) contracts.RouteGroup {
			basePath := usersServer.Cfg().GetBasePath()
			v1 := usersServer.RouteBuilder().Group(basePath)
			return v1.Group("/auth")
		}, fx.ResultTags(`name:"auth-routes"`))),

	// Auth endpoints
	fx.Provide(
		route.AsRoute(auth.NewSignUpHandler, "auth-routes"),
		route.AsRoute(auth.NewSignInHandler, "auth-routes"),
		route.AsRoute(auth.NewChangePasswordHandler, "auth-routes"),
	),

	// User/Profile endpoints
	fx.Provide(
		route.AsRoute(profile.NewGetUserHandler, "user-routes"),
		route.AsRoute(profile.NewUpdateProfileHandler, "user-routes"),
		route.AsRoute(profile.NewDisableUserHandler, "user-routes"),
		route.AsRoute(profile.NewEnableUserHandler, "user-routes"),
	),
)
