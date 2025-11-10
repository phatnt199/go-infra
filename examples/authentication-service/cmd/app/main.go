package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/app"
	_ "github.com/phatnt199/go-infra/pkg/component/authentication/handlers"
)

// @title Authentication Service API
// @version 1.0
// @description Authentication service using identifier/credential pattern
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	application := app.NewApplication()
	application.Run()
}
