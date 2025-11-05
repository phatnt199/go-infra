package v1

import (
	"net/http"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1/fxparams"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/core/web/route"
)

type createUserHandler struct {
	fxparams.UserRouteParams
}

func NewCreateUserHandler(
	params fxparams.UserRouteParams,
) route.Endpoint {
	return &createUserHandler{
		UserRouteParams: params,
	}
}

func (ep *createUserHandler) MapEndpoint() {
	ep.UsersGroup.POST("", ep.handler())
}

func (ep *createUserHandler) handler() contracts.HandlerFunc {
	return func(c contracts.Context) error {
		ctx := c.Request().Context()

		return c.JSON(http.StatusCreated, ctx.Done())
	}
}
