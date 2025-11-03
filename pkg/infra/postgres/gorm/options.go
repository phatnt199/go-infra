package postgresgorm

import (
	"github.com/iancoleman/strcase"
	"github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	typeMapper "github.com/phatnt199/go-infra/pkg/reflection/typemapper"
)

type GormType int

const (
	Postgres GormType = iota
	SQLite
	InMemory
)

type GormOptions struct {
	Type          GormType `mapstructure:"type"`
	Host          string   `mapstructure:"host"`
	Port          int      `mapstructure:"port"`
	User          string   `mapstructure:"user"`
	Password      string   `mapstructure:"password"`
	DBName        string   `mapstructure:"dbname"`
	SSLMode       bool     `mapstructure:"sslmode"`
	EnableTracing bool     `mapstructure:"enable_tracing"`
}

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[GormOptions]())

func provideConfig(environment environment.Environment) (*GormOptions, error) {
	return config.BindConfigKey[*GormOptions](optionName, environment)
}
