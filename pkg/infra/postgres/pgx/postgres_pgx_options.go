package pgx

import (
	"github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	typeMapper "github.com/phatnt199/go-infra/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[PostgresPgxOptions]())

type PostgresPgxOptions struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  bool   `mapstructure:"sslMode"`
	Password string `mapstructure:"password"`
	LogLevel int    `mapstructure:"logLevel"`
}

func provideConfig(environment environment.Environment) (*PostgresPgxOptions, error) {
	return config.BindConfigKey[*PostgresPgxOptions](optionName, environment)
}
