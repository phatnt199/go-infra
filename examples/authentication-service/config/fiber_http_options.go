package config

import (
	"strings"

	"github.com/iancoleman/strcase"
	appConfig "github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	"github.com/phatnt199/go-infra/pkg/reflection/typemapper"
)

type FiberHttpOptions struct {
	Host        string `mapstructure:"host" env:"HTTP_HOST" default:"0.0.0.0"`
	Port        string `mapstructure:"port" env:"HTTP_PORT" default:":8080"`
	BasePath    string `mapstructure:"basePath" env:"HTTP_BASE_PATH" default:"/api/v1"`
	Name        string `mapstructure:"name" env:"SERVICE_NAME" default:"Authentication Service"`
	Development bool   `mapstructure:"development" env:"DEVELOPMENT" default:"true"`
	Timeout     int    `mapstructure:"timeout" env:"HTTP_TIMEOUT" default:"30"`
}

func NewFiberHttpOptions(environment environment.Environment) (*FiberHttpOptions, error) {
	optionName := strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[FiberHttpOptions]())

	cfg, err := appConfig.BindConfigKey[*FiberHttpOptions](optionName, environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (o *FiberHttpOptions) String() string {
	return strings.Join([]string{o.Host, o.Port}, ":")
}
