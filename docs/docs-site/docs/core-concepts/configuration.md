---
sidebar_position: 4
---

# Configuration Management

Learn how to manage application configuration in go-infra using environment-based settings.

## Overview

go-infra provides a flexible configuration system that supports:

- **Environment-based configs** - Different settings for dev, staging, production
- **Multiple sources** - Files, environment variables, command-line flags
- **Type-safe** - Strongly typed configuration structs
- **Validation** - Built-in validation support
- **Hot reload** - Update config without restart (optional)

## Quick Start

### Define Configuration Structure

```go
// config/config.go
package config

type Config struct {
    Environment string         `json:"environment" env:"APP_ENV" default:"development"`
    Server      ServerConfig   `json:"server"`
    Database    DatabaseConfig `json:"database"`
    JWT         JWTConfig      `json:"jwt"`
    Email       EmailConfig    `json:"email"`
}

type ServerConfig struct {
    Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
    Port int    `json:"port" env:"SERVER_PORT" default:"3000"`
}

type DatabaseConfig struct {
    Host     string `json:"host" env:"DB_HOST" default:"localhost"`
    Port     int    `json:"port" env:"DB_PORT" default:"5432"`
}
```

## Bootstrapping the environment

Go-infra provides two complementary configuration experiences:

1. **Bootstrapping the environment** via `pkg/application/environment`, which loads `.env`, resolves `APP_ENV`, sets `APP_ROOT_PATH`, and exposes helpers for working with process variables.
2. **Typed configuration structs** under `pkg/application/config`, which either bind a `config.<env>` file or read environment variables directly, complete with defaults, helper methods, and validation.

Use whichever flow matches your service: many examples combine both (`ConfigAppEnv()` to prepare the environment and `config.BindConfig*` to load structured files), but you can just consume `config.Load` if all settings come from environment variables.

## Key capabilities

- **Environment discovery:** `environment.ConfigAppEnv()` recursively loads `.env` files using `godotenv`, writes the resolved `APP_ENV`/`APP_NAME`/`APP_ROOT_PATH` into Viper, and optionally changes the working directory through `FixProjectRootWorkingDirectoryPath()`.
- **File binding with Viper:** `config.BindConfig`/`BindConfigKey` sets defaults via `go-defaults`, picks the right `config.<env>` file, unmarshals it through Viper, runs `viper.AutomaticEnv()`, and finally `env.Parse(cfg)` to honor `env` struct tags.
- **Environment-only loader:** `config.Load()` (and helpers `LoadOnce`, `Get`, `Set`) reads environment variables (`APP_*`, `HTTP_*`, `DB_*`, `LOG_*`, `JWT_*`, etc.), applies defaults, validates each section (`validation.go`), and exposes helper methods such as `cfg.Server.HTTP.Address()` and `cfg.Database.DSN()`.

### Example: Bootstrapping

1. Start your application by calling `environment.ConfigAppEnv()`. It:
   - loads a `.env` file located inside the current working directory or any parent directory,
   - logs the detected `APP_ENV` and falls back to `development` when nothing is set,
   - sets `APP_ROOT_PATH` based on either `APP_NAME` (if provided) or by searching for the nearest `go.mod`, and
   - calls `FixProjectRootWorkingDirectoryPath()` so later helpers can rely on a stable root.
2. Pass the returned `environment.Environment` value to other packages (or store it globally) so they all share the same view of `development`, `staging`, or `production`.
3. Optionally set `CONFIG_PATH` to the directory that holds your `config.<env>.json` file. Without it, the binder walks from `APP_ROOT_PATH` until it finds the first directory that contains any `config.<env>.(json|yaml|yml)` file.
4. Remember: `APP_ENV` in the OS environment overrides explicit values passed into the binder, so unset it only when you want the code-provided value to remain authoritative.

```go
env := environment.ConfigAppEnv()
cfg, err := config.BindConfigKey[*AppOptions](optionName, env)
if err != nil {
    log.Fatalf("unable to bind config: %v", err)
}
```

## Binding configuration files

`pkg/application/config/config_helper.go` handles file-based configuration.

1. `BindConfig`/`BindConfigKey` creates the zero value of your struct and sets defaults through `go-defaults`.
2. It determines the environment (explicit parameter or `constants.DEV_ENV`, then `APP_ENV` overrides) and resolves the directory where `config.<env>` lives:
   - use `CONFIG_PATH` to point directly to the folder containing `config.development.json` (recommended for CI/Test), or
   - omit `CONFIG_PATH` so the helper searches from `APP_ROOT_PATH` downwards for the first matching file.
3. Viper loads `config.<env>.json` (the type is currently `json`) inside that directory, unmarshals it into your struct, and, after the file read, runs `viper.AutomaticEnv()`.
4. After Viper, `env.Parse(cfg)` reads environment variables according to your `env:"..."` tags and overwrites any fields accordingly.
5. Pass a `configKey` to read a named section instead of the entire file (examples use the `typemapper`/`strcase` helpers to derive keys such as `fiberHttpOptions`).

```go
type AppOptions struct {
    ServiceName  string `mapstructure:"serviceName" env:"ServiceName"`
    DeliveryType string `mapstructure:"deliveryType" env:"DeliveryType" default:"http"`
}

func NewAppOptions(env environment.Environment) (*AppOptions, error) {
    optionName := strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[AppOptions]())
    return config.BindConfigKey[*AppOptions](optionName, env)
}
```

Any environment variable that matches the `env` tag (case sensitive) or a Viper key wins over the file contents, and the helper prints parsing errors to help you diagnose missing overrides.

## Environment-driven loader (`pkg/application/config`)

If your service prefers to configure itself entirely from environment variables, `config.Load()` is the standard entry point. It builds a `Config` with the following sections (all validated via `validation.go` & defaults baked into the loader functions):

- App — `APP_NAME`, `APP_VERSION`, `APP_ENV`, `APP_DEBUG`, `APP_TIMEZONE`, `APP_TIMEOUT`
- Server — HTTP (host/port/timeouts/CORS/TLS) and gRPC (host/port/keepalives)
- Database — driver, host, port, credentials, pooling, migration path
- Redis — host, port, password, pool settings, TLS
- Queue — driver, URL, concurrency, retries
- Storage — driver (S3/MinIO/GCS/local), credentials, bucket, SSL
- Logger — log level/format/output paths, caller/stacktrace toggles that adjust defaults based on `APP_ENV`
- Auth — JWT secrets/keys, OAuth providers, session cookies, password policies

Helper methods such as `cfg.App.IsProduction()`, `cfg.Server.HTTP.Address()`, `cfg.Database.DSN()`, and `cfg.Redis.Address()` simplify downstream code.

```go
cfg, err := config.Load()
if err != nil {
    log.Fatalf("config error: %v", err)
}
log.Printf("starting %s/%s", cfg.App.Name, cfg.App.Environment)
httpServer := http.Server{
    Addr:         cfg.Server.HTTP.Address(),
    ReadTimeout:  cfg.Server.HTTP.ReadTimeout,
    WriteTimeout: cfg.Server.HTTP.WriteTimeout,
}
```

Use `config.LoadOnce()`/`config.Get()` when you need a cached version (e.g., tests). Call `config.Set` with a handcrafted `*config.Config` to inject deterministic values during unit tests.

## Example workflow

1. Copy or create a `.env` file in your repo or export the relevant environment variables (APP*ENV, HTTP*_, DB\__, etc.).
2. Run `go run ./examples/env-usage` (or your service entrypoint) — the example illustrates four modes (default, explicit env, `.env`/APP_ENV overrides, and config binding).
3. When a config file is involved, place it beside your service config directory (see `examples/users-api/internal/config/config.development.json`).
4. Prefer setting `CONFIG_PATH` before running the binder (`export CONFIG_PATH=$(pwd)/examples/users-api/internal/config`) so you avoid the repo scan, and then call `config.BindConfig*`.
5. The builder layers expect you to reuse the same `environment.Environment` value so logging, metrics, and other components can read the same `APP_ENV` value.

## Best practices & troubleshooting

- **Pin `CONFIG_PATH` or pass explicit directories in CLI/test scripts** so Viper doesn’t perform a costly filesystem walk and you always load the right file.
- **Remember `FixProjectRootWorkingDirectoryPath()` can `chdir`**; running from the repo root and setting `CONFIG_PATH` avoids surprises with `.env` loading.
- **Keep secrets out of source control** and prefer environment variables, `env` tags, or injected secrets (Vault/SecretsManager) when necessary.
- **Override after file loads:** config files are read first, followed by `viper.AutomaticEnv()` and `env.Parse`, so the most specific override wins.
- **Validation errors are helpful:** the loader aggregates field-level messages (e.g., `server.http.port` or `auth.jwt.secret`). Treat them as fatal during startup.
- **Use the helpers:** `cfg.Server.HTTP.Address()` and `cfg.Database.DSN()` avoid repeated formatting logic.

## Related examples

- `examples/env-usage` documents the `.env` discovery logic, how `APP_ENV` interacts with explicit parameters, and how to force `CONFIG_PATH` for reproducible runs.
- `examples/authentication-service/config` and `examples/microservices/userservice/config` show real config structs (with `mapstructure`/`env` tags) and how to call `config.BindConfigKey` for each section.

## Learn more

- The rest of the docs site explores deployment, observability, and error handling; when you add a new config option, add it under `pkg/application/config` and ensure the validator surface is updated.
