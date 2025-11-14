---
sidebar_position: 3
---

# Dependency Injection

go-infra builds every application with [Uber Fx](https://github.com/uber-go/fx) so that the container automatically instantiates, wires, and tears down dependencies. This page describes how to use the `fxapp` builder and the Fx primitives that ship with go-infra instead of Google Wire or hand-rolled containers.

## go-infra Fx workflow

Every go-infra entry point arranges its dependencies through `github.com/phatnt199/go-infra/pkg/adapter/fxapp`. The library sets up a logger, reads environment configuration, and builds an `fx.App` under the hood. Your job is to offer modules, constructors, and hooks that plug into that container:

1. Create an `fxapp.ApplicationBuilder` via `fxapp.NewApplicationBuilder(...)`. It captures the desired environment and creates a default logger and config module.
2. Provide application-wide functionality using `ProvideModule` (for reusable `fx.Module`s), `Provide` (for raw constructors), `Decorate`, or custom parameterized hooks.
3. Call `Build()` to get a `contracts.Application`, register any extra lifecycle hooks with `RegisterHook`, and then `Run()` the application.

```go
appBuilder := fxapp.NewApplicationBuilder()
appBuilder.ProvideModule(config.Module)
appBuilder.ProvideModule(customfiber.Module)
appBuilder.ProvideModule(postgresgorm.Module)
appBuilder.ProvideModule(modules.Module)

app := appBuilder.Build()
app.RegisterHook(setupSwagger)
app.Run()
```

`modules.Module` in this example comes from `examples/users-api/internal/modules` and demonstrates the recommended pattern for defining go-infra modules.

## Defining fx Modules

An Fx module groups related constructors, decorators, and invocations. go-infra modules usually follow this shape:

```go
var Module = fx.Module(
    "users_api_module",
    fx.Provide(repository.NewUserRepository),
    fx.Provide(handler.NewUserHandler),
    fx.Invoke(setupRoutes),
)
```

The module registers the repository, handler, and route wiring inside a single `fx.Module`. The `setupRoutes` function later consumes concrete dependencies:

```go
func setupRoutes(server contracts.HttpServer, userHandler *handler.UserHandler, logger logger.Logger) {
    routes.SetupRoutes(server, userHandler, logger)
}
```

Combine modules (see `examples/microservices/userservice/internal/shared/configurations/users/users_fx.go`) by nesting them inside a parent module.

### Providers with fx.In / fx.Out

go-infra heavily uses `fx.In` and `fx.Out` structs so providers can express what they consume and produce. For example, the authentication module returns several services:

```go
type ProvideAuthenticationOptions struct {
    fx.In

    DBContext  dbcontext.AuthGormDBContext
    AuthConfig *config.AuthOptions
    Logger     logger.Logger
}

type ProvideAuthenticationResult struct {
    fx.Out

    AuthComponent *authComponent.Component
    AuthService   authContracts.IAuthService
    TokenService  authContracts.ITokenService
    UserProvider  authContracts.IUserProvider
}

func provideAuthentication(opts ProvideAuthenticationOptions) ProvideAuthenticationResult {
    authComp := authComponent.NewComponentWithConfig(...)
    return ProvideAuthenticationResult{
        AuthComponent: authComp,
        AuthService:   authComp.GetAuthService(),
        TokenService:  authComp.GetTokenService(),
        UserProvider:  provider.NewUserProvider(opts.DBContext),
    }
}
```

Use `fx.Annotate` together with `fx.ParamTags` when you need to resolve named bindings or multiple implementations of the same interface. The builder exposes `ResolveFuncWithParamTag` so you can attach handlers that rely on those tags.

## Hooks, decorators, and lifecycle

Go-infra exposes the Fx container through `contracts.Application`, which allows you to:

- `ResolveFunc` / `ResolveFuncWithParamTag` to run Fx functions at build time (useful for wiring routers or metrics)
- `RegisterHook` to add lifecycle hooks that depend on container-provided services (e.g., registering Swagger or health checks).

The HTTP setup hook in `examples/users-api/cmd/api/main.go` illustrates how to read the Fiber app instance from the `contracts.HttpServer` and add Swagger routes while logging to the shared logger.

## Scoped dependencies and testing

For request-scoped values, materialize them via middleware that resolves dependencies from the container or Fx parameters. You can also decorate constructors to wrap or enhance dependencies before the application starts.

Testing stays focused on constructors by instantiating them directly with mocks. Because interfaces are the default dependency, you can swap real implementations for fakes in unit tests, even though Fx wires the production application.

## Best practices

- **Prefer modules**: group related providers and invocations into `fx.Module` so you can reuse them across services.
- **Return interfaces** from constructors and avoid leaking concrete types. The go-infra packages expose contracts (e.g., `contracts.HttpServer`, `authContracts.IAuthService`) for this reason.
- **Keep constructors small**: don’t open database connections or mutate global state inside constructors; rely on Fx lifecycle hooks when needed.
- **Use the builder**: the `fxapp.ApplicationBuilder` already wires configuration (`config.Module`) and logging (`zap.ModuleFunc`) so you can focus on feature modules.

## Resources

- `examples/users-api/cmd/api/main.go`: Application entry point that assembles modules, hooks, and Swagger.
- `examples/authentication-service/internal/auth/auth_fx.go`: Demonstrates `fx.In`/`fx.Out` constructors inside a feature component.
- `pkg/adapter/fxapp`: Library that bridges go-infra primitives (`config`, `logger`, `environment`) with Fx.---
  sidebar_position: 3

---

# Dependency Injection

Learn how to use dependency injection patterns in go-infra applications.

## Why Dependency Injection?

Dependency Injection (DI) is a design pattern that helps you:

- **Decouple components** - Reduce tight coupling between modules
- **Improve testability** - Easy to mock dependencies in tests
- **Enhance flexibility** - Swap implementations without changing code
- **Better organization** - Clear dependency hierarchy

## Constructor Injection

        logger: logger,
    // Initialize repositories

---

## sidebar_position: 3

# Dependency Injection

go-infra builds every application with [Uber Fx](https://github.com/uber-go/fx) so that the container automatically instantiates, wires, and tears down dependencies. This page describes how to use the `fxapp` builder and the Fx primitives that ship with go-infra instead of Google Wire or hand-rolled containers.

## go-infra Fx workflow

Every go-infra entry point arranges its dependencies through `github.com/phatnt199/go-infra/pkg/adapter/fxapp`. The library sets up a logger, reads environment configuration, and builds an `fx.App` under the hood. Your job is to offer modules, constructors, and hooks that plug into that container:

1. Create an `fxapp.ApplicationBuilder` via `fxapp.NewApplicationBuilder(...)`. It captures the desired environment and creates a default logger and config module.
2. Provide application-wide functionality using `ProvideModule` (for reusable `fx.Module`s), `Provide` (for raw constructors), `Decorate`, or custom parameterized hooks.
3. Call `Build()` to get a `contracts.Application`, register any extra lifecycle hooks with `RegisterHook`, and then `Run()` the application.

```go
appBuilder := fxapp.NewApplicationBuilder()
appBuilder.ProvideModule(config.Module)
appBuilder.ProvideModule(customfiber.Module)
appBuilder.ProvideModule(postgresgorm.Module)
appBuilder.ProvideModule(modules.Module)

app := appBuilder.Build()
app.RegisterHook(setupSwagger)
app.Run()
```

`modules.Module` in this example comes from `examples/users-api/internal/modules` and demonstrates the recommended pattern for defining go-infra modules.

## Defining fx Modules

An Fx module groups related constructors, decorators, and invocations. go-infra modules usually follow this shape:

```go
var Module = fx.Module(
    "users_api_module",
    fx.Provide(repository.NewUserRepository),
    fx.Provide(handler.NewUserHandler),
    fx.Invoke(setupRoutes),
)
```

The module registers the repository, handler, and route wiring inside a single `fx.Module`. The `setupRoutes` function later consumes concrete dependencies:

```go
func setupRoutes(server contracts.HttpServer, userHandler *handler.UserHandler, logger logger.Logger) {
    routes.SetupRoutes(server, userHandler, logger)
}
```

Combine modules (see `examples/microservices/userservice/internal/shared/configurations/users/users_fx.go`) by nesting them inside a parent module.

### Providers with fx.In / fx.Out

go-infra heavily uses `fx.In` and `fx.Out` structs so providers can express what they consume and produce. For example, the authentication module returns several services:

```go
type ProvideAuthenticationOptions struct {
    fx.In

    DBContext  dbcontext.AuthGormDBContext
    AuthConfig *config.AuthOptions
    Logger     logger.Logger
}

type ProvideAuthenticationResult struct {
    fx.Out

    AuthComponent *authComponent.Component
    AuthService   authContracts.IAuthService
    TokenService  authContracts.ITokenService
    UserProvider  authContracts.IUserProvider
}

func provideAuthentication(opts ProvideAuthenticationOptions) ProvideAuthenticationResult {
    authComp := authComponent.NewComponentWithConfig(...)
    return ProvideAuthenticationResult{
        AuthComponent: authComp,
        AuthService:   authComp.GetAuthService(),
        TokenService:  authComp.GetTokenService(),
        UserProvider:  provider.NewUserProvider(opts.DBContext),
    }
}
```

Use `fx.Annotate` together with `fx.ParamTags` when you need to resolve named bindings or multiple implementations of the same interface. The builder exposes `ResolveFuncWithParamTag` so you can attach handlers that rely on those tags.

## Hooks, decorators, and lifecycle

Go-infra exposes the Fx container through `contracts.Application`, which allows you to:

- `ResolveFunc` / `ResolveFuncWithParamTag` to run Fx functions at build time (useful for wiring routers or metrics)
- `RegisterHook` to add lifecycle hooks that depend on container-provided services (e.g., registering Swagger or health checks).

The HTTP setup hook in `examples/users-api/cmd/api/main.go` illustrates how to read the Fiber app instance from the `contracts.HttpServer` and add Swagger routes while logging to the shared logger.

## Scoped dependencies and testing

For request-scoped values, materialize them via middleware that resolves dependencies from the container or Fx parameters. You can also decorate constructors to wrap or enhance dependencies before the application starts.

Testing stays focused on constructors by instantiating them directly with mocks. Because interfaces are the default dependency, you can swap real implementations for fakes in unit tests, even though Fx wires the production application.

## Best practices

- **Prefer modules**: group related providers and invocations into `fx.Module` so you can reuse them across services.
- **Return interfaces** from constructors and avoid leaking concrete types. The go-infra packages expose contracts (e.g., `contracts.HttpServer`, `authContracts.IAuthService`) for this reason.
- **Keep constructors small**: don’t open database connections or mutate global state inside constructors; rely on Fx lifecycle hooks when needed.
- **Use the builder**: the `fxapp.ApplicationBuilder` already wires configuration (`config.Module`) and logging (`zap.ModuleFunc`) so you can focus on feature modules.

## Resources

- `examples/users-api/cmd/api/main.go`: Application entry point that assembles modules, hooks, and Swagger.
- `examples/authentication-service/internal/auth/auth_fx.go`: Demonstrates `fx.In`/`fx.Out` constructors inside a feature component.
- `pkg/adapter/fxapp`: Library that bridges go-infra primitives (`config`, `logger`, `environment`) with Fx.
