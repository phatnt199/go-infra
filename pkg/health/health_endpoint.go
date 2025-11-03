package health

import (
	"net/http"

	httpContracts "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	healthContracts "github.com/phatnt199/go-infra/pkg/health/contracts"
)

type HealthCheckEndpoint struct {
	service    healthContracts.HealthService
	httpServer httpContracts.HttpServer
}

func NewHealthCheckEndpoint(
	service healthContracts.HealthService,
	server httpContracts.HttpServer,
) *HealthCheckEndpoint {
	return &HealthCheckEndpoint{service: service, httpServer: server}
}

func (s *HealthCheckEndpoint) RegisterEndpoints() {
	s.httpServer.RouteBuilder().GET("health", s.checkHealth)
}

func (s *HealthCheckEndpoint) checkHealth(c httpContracts.Context) error {
	check := s.service.CheckHealth(c.Request().Context())
	if !check.AllUp() {
		return c.JSON(http.StatusServiceUnavailable, check)
	}

	return c.JSON(http.StatusOK, check)
}
