# env-usage example

This example demonstrates how the repository resolves environment and configuration files. It includes four runnable examples (A–D) showing common usage patterns:

- Example A — Default: call `fxapp.NewApplicationBuilder()` with no arguments. This will default to `development` unless `APP_ENV` is set in the process environment or in a loaded `.env` file.
- Example B — Explicit: call `fxapp.NewApplicationBuilder(environment.Production)` to force the environment from code. Note the current implementation still allows the process `APP_ENV` to override this value (see "Precedence" below).
- Example C — APP_ENV/.env: demonstrate reading `APP_ENV` from the operating system environment or from a `.env` file that the code loads at startup.
- Example D — Config file binding: demonstrates loading `config.development.json` (the Users API example config) via the repository's config binding helpers and printing the fiber HTTP options.

This README explains how to use `.env`, how config files are located, how to override values, and common troubleshooting tips.

## How config files are found

The config binding logic (see `pkg/application/config/config_helper.go`) uses the following rules:

1. Determine the environment name to load (`development`, `staging`, `production`):

   - If you pass an explicit `environment.Environment` to the helper, it will be used initially.
   - The code also reads the `APP_ENV` process environment variable and will override the explicit value if present (see "Precedence" below).

2. Determine the config directory (`configPath`):

   - If `CONFIG_PATH` (viper key `ConfigPath`) is set, the loader will use that directory directly.
   - Otherwise it uses `APP_ROOT_PATH` (set by `ConfigAppEnv`) or searches up from the project root (the helper searches for the first directory containing `config.<env>.json` under that root).

3. Viper loads a file named `config.<env>` (for example `config.development.json`) from the chosen `configPath` and unmarshals it into the target struct.

4. After file values are loaded, the code calls `viper.AutomaticEnv()` and `env.Parse(cfg)` so environment variables can override struct fields (fields should have `env` tags if you want `env.Parse` to map them).

## Running the example

From the repository root you can run the example program:

```sh
go run ./examples/env-usage
```

### Run with `.env` (local development)

1. Copy the example `.env` to the repository root (or set `APP_ENV` in your shell):

```sh
cp ./examples/env-usage/.env.example ./.env
# or
export APP_ENV=production
```

2. Run:

```sh
go run ./examples/env-usage
```

You should see Example A/C logs showing `APP_ENV` and Example D printing fiber options loaded from `config.development.json` (or the corresponding file for the environment you chose).

### Force a config directory (fast and deterministic)

If you want to avoid the repository-wide search, set `CONFIG_PATH` to the directory that contains `config.<env>.json`:

```sh
export CONFIG_PATH=$(pwd)/examples/users-api/internal/config
go run ./examples/env-usage
```

This is the recommended approach for CI or reproducible local runs.

## Precedence and important caveats

- The loader reads file values first, then applies environment overrides.
- Currently `ConfigAppEnv` will override an explicit environment parameter if `APP_ENV` is present in the process environment. This behavior can be surprising — if you want the explicit parameter to be authoritative, call `fxapp.NewApplicationBuilder(environment.X)` and ensure `APP_ENV` is not set, or request a change to the code to make the explicit parameter final.
- The code also calls `FixProjectRootWorkingDirectoryPath()` which calls `os.Chdir(rootDir)` and therefore changes the process working directory. Because `.env` loading searches from the current working directory, this can change whether subsequent `.env` files are found. To avoid this class of surprise:
  - Prefer setting `CONFIG_PATH` explicitly.
  - Or always run from the repository root.

## How to override individual values with environment variables

- After the config file is loaded into a struct, `env.Parse(cfg)` runs and will map environment variables to struct fields based on `env` tags (if present). For example, if a config struct field has `env:"Port"`, set `PORT` (or the exact mapping used) to override it.
- You can also set typical config environment variables used by the app loader, such as `APP_ENV`, `APP_NAME`, `CONFIG_PATH` and logger-related env vars.

## Troubleshooting

- If you see `viper.ReadInConfig` errors, check `CONFIG_PATH` and make sure `config.<env>.json` exists there.
- If `.env` appears not to be loaded on subsequent runs, remember the loader may have moved the working directory; check the initial working directory and whether `APP_ENV` remains set in the process environment.
- To debug, try:
  ```sh
  # print process APP_ENV and explicit config path
  echo APP_ENV=$APP_ENV
  go run ./examples/env-usage
  ```

## Tests and CI

- For tests prefer passing an explicit environment to the builder or set `CONFIG_PATH` so tests are isolated and deterministic.
- Avoid relying on `os.Chdir` side-effects in libraries — if you write tests that depend on the working directory, run them from a consistent path or change the code to avoid `os.Chdir`.

## Example commands summary

Run example normally (use .env if present in repo root):

```sh
go run ./examples/env-usage
```

Force config path:

```sh
export CONFIG_PATH=$(pwd)/examples/users-api/internal/config
go run ./examples/env-usage
```

Force APP_ENV in the shell (overrides explicit param in current implementation):

```sh
export APP_ENV=production
go run ./examples/env-usage
```
