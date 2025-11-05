package users

import (
	"github.com/phatnt199/go-infra/examples/microservices/userservice/config"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/configurations/users/infrastructure"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/data"
	"go.uber.org/fx"
)

var UsersServiceModule = fx.Module(
	"usersfx",

	config.Module,
	infrastructure.Module,
	data.Module,
)
