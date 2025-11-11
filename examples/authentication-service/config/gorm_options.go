package config

import (
	"strings"

	"github.com/iancoleman/strcase"
	appConfig "github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	"github.com/phatnt199/go-infra/pkg/reflection/typemapper"
)

type GormOptions struct {
	Type          int    `mapstructure:"type" env:"GORM_TYPE" default:"0"`
	Host          string `mapstructure:"host" env:"GORM_HOST" default:"localhost"`
	Port          int    `mapstructure:"port" env:"GORM_PORT" default:"54100"`
	User          string `mapstructure:"user" env:"GORM_USER" default:"admin"`
	Password      string `mapstructure:"password" env:"GORM_PASSWORD" default:"123456"`
	DbName        string `mapstructure:"dbName" env:"GORM_DBNAME" default:"auth_db"`
	SSLMode       bool   `mapstructure:"sslmode" env:"GORM_SSLMODE" default:"false"`
	EnableTracing bool   `mapstructure:"enableTracing" env:"GORM_ENABLE_TRACING" default:"false"`
}

func NewGormOptions(environment environment.Environment) (*GormOptions, error) {
	optionName := strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[GormOptions]())

	cfg, err := appConfig.BindConfigKey[*GormOptions](optionName, environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (o *GormOptions) String() string {
	password := o.Password
	if password != "" {
		password = "***"
	}
	return strings.Join([]string{o.User, password, "@", o.Host, ":", string(rune(o.Port))}, "")
}
