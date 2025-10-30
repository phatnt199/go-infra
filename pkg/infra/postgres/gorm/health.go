package postgresgorm

import (
	"context"
	"database/sql"

	"github.com/phatnt199/go-infra/pkg/health/contracts"
)

type gormHealthChecker struct {
	client *sql.DB
}

func NewGormHealthChecker(client *sql.DB) contracts.Health {
	return &gormHealthChecker{client}
}

func (healthChecker *gormHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.PingContext(ctx)
}

func (healthChecker *gormHealthChecker) GetHealthName() string {
	return "postgres"
}
