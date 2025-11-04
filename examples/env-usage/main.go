package main

import (
	"fmt"
	"os"

	fxapp "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
	fibercfg "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
)

func main() {
	fmt.Println("env-usage example: demonstrating how ConfigAppEnv and fxapp.NewApplicationBuilder resolve environment")

	// Example A: default behavior (no args, no APP_ENV set)
	fmt.Println("\n--- Example A: default (no args passed to NewApplicationBuilder) ---")
	a := fxapp.NewApplicationBuilder()
	loggerA := a.Logger()
	envA := a.Environment()
	loggerA.Infof("Resolved environment: %s", envA.GetEnvironmentName())
	fmt.Printf("IsDevelopment=%v, IsProduction=%v, IsStaging=%v\n", envA.IsDevelopment(), envA.IsProduction(), envA.IsStaging())

	// Example B: explicit environment passed to NewApplicationBuilder
	fmt.Println("\n--- Example B: explicit environment (pass environment.Production) ---")
	b := fxapp.NewApplicationBuilder(environment.Production)
	loggerB := b.Logger()
	envB := b.Environment()
	loggerB.Infof("Resolved environment: %s", envB.GetEnvironmentName())
	fmt.Printf("IsDevelopment=%v, IsProduction=%v, IsStaging=%v\n", envB.IsDevelopment(), envB.IsProduction(), envB.IsStaging())

	// Example C: demonstrate reading APP_ENV from the OS environment
	fmt.Println("\n--- Example C: APP_ENV from environment (set APP_ENV or use .env file) ---")
	// Print current APP_ENV seen by the process (this may be set via OS env or by a loaded .env file)
	fmt.Printf("OS APP_ENV value: %q\n", os.Getenv("APP_ENV"))

	c := fxapp.NewApplicationBuilder()
	loggerC := c.Logger()
	envC := c.Environment()
	loggerC.Infof("Resolved environment: %s", envC.GetEnvironmentName())
	fmt.Printf("IsDevelopment=%v, IsProduction=%v, IsStaging=%v\n", envC.IsDevelopment(), envC.IsProduction(), envC.IsStaging())

	fmt.Println("\nNotes: see README.md in this folder for instructions on using .env or explicit env values.")

	// Example D: load users-api fiber HTTP options from example project's config.development.json
	fmt.Println("\n--- Example D: load Users API Fiber options from config.development.json ---")
	// We explicitly request the Development environment so the loader searches for config.development.json
	fiberOptions, err := fibercfg.ProvideConfig(environment.Development)
	if err != nil {
		fmt.Printf("failed to load fiber config: %v\n", err)
	} else {
		fmt.Printf("Fiber Port: %s, Host: %s, BasePath: %s, Name: %s, Development: %v\n",
			fiberOptions.Port, fiberOptions.Host, fiberOptions.BasePath, fiberOptions.Name, fiberOptions.Development)
	}
}
