package customfiber

import (
	"context"
	"errors"
	"net/http"

	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter/config"
	"github.com/phatnt199/go-infra/pkg/logger"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx"
)

// FiberHttpServerParams contains the parameters for NewFiberHttpServer with optional meter
type FiberHttpServerParams struct {
	fx.In

	Cfg    *config.FiberHttpOptions
	Logger logger.Logger
	Meter  metric.Meter `optional:"true"`
}

// ProvideFiberServer creates the HTTP server with fallback to noop meter if not provided
func ProvideFiberServer(params FiberHttpServerParams) contracts.HttpServer {
	// Use noop meter if not provided (when metrics module is not loaded)
	if params.Meter == nil {
		params.Meter = noop.MeterProvider{}.Meter("")
	}
	return NewFiberHttpServer(params.Cfg, params.Logger, params.Meter)
}

var (
	// Module provides Fiber HTTP server using fx dependency injection
	Module = fx.Module(
		"fiberfx",
		fiberProviders,
		fiberInvokes,
	)

	fiberProviders = fx.Options(
		fx.Provide(
			config.ProvideConfig,
			fx.Annotate(
				ProvideFiberServer,
				fx.As(new(contracts.HttpServer)),
			),
		),
	)

	fiberInvokes = fx.Options(fx.Invoke(registerHooks))
)

// registerHooks registers lifecycle hooks for the Fiber server
func registerHooks(
	lc fx.Lifecycle,
	fiberServer contracts.HttpServer,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := fiberServer.RunHttpServer(); !errors.Is(err, http.ErrServerClosed) {
					logger.Fatalf(
						"(FiberHttpServer.RunHttpServer) error in running server: {%v}",
						err,
					)
				}
			}()
			fiberServer.Logger().Infof(
				"%s is listening on Host:{%s} Http PORT: {%s}",
				fiberServer.Cfg().GetName(),
				fiberServer.Cfg().GetHost(),
				fiberServer.Cfg().GetPort(),
			)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := fiberServer.GracefulShutdown(ctx); err != nil {
				fiberServer.Logger().Errorf("error shutting down fiber server: %v", err)
			} else {
				fiberServer.Logger().Info("fiber server shutdown gracefully")
			}
			return nil
		},
	})
}
