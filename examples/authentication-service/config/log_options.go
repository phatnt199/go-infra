package config

import (
	"strings"

	"github.com/iancoleman/strcase"
	appConfig "github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	"github.com/phatnt199/go-infra/pkg/reflection/typemapper"
)

type LogOptions struct {
	Level         string `mapstructure:"level" env:"LOG_LEVEL" default:"info"`
	LogType       int    `mapstructure:"logType" env:"LOG_TYPE" default:"0"`
	CallerEnabled bool   `mapstructure:"callerEnabled" env:"LOG_CALLER_ENABLED" default:"true"`
	EnableTracing bool   `mapstructure:"enableTracing" env:"LOG_ENABLE_TRACING" default:"false"`
}

func NewLogOptions(environment environment.Environment) (*LogOptions, error) {
	optionName := strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[LogOptions]())

	cfg, err := appConfig.BindConfigKey[*LogOptions](optionName, environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (o *LogOptions) String() string {
	return strings.ToUpper(o.Level)
}
