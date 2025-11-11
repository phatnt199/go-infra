package config

import (
	"strings"

	"github.com/iancoleman/strcase"
	appConfig "github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	"github.com/phatnt199/go-infra/pkg/reflection/typemapper"
)

type AuthOptions struct {
	JWT JWT `mapstructure:"jwt"`
}

type JWT struct {
	Secret        string `mapstructure:"secret" env:"JWT_SECRET" default:"your-secret-key-change-in-production"`
	Issuer        string `mapstructure:"issuer" env:"JWT_ISSUER" default:"authentication-service"`
	Audience      string `mapstructure:"audience" env:"JWT_AUDIENCE" default:"authentication-api"`
	Algorithm     string `mapstructure:"algorithm" env:"JWT_ALGORITHM" default:"HS256"`
	AccessExpiry  string `mapstructure:"accessExpiry" env:"JWT_ACCESS_EXPIRY" default:"15m"`
	RefreshExpiry string `mapstructure:"refreshExpiry" env:"JWT_REFRESH_EXPIRY" default:"168h"`
}

func NewAuthOptions(environment environment.Environment) (*AuthOptions, error) {
	optionName := strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[AuthOptions]())

	cfg, err := appConfig.BindConfigKey[*AuthOptions](optionName, environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (o *AuthOptions) String() string {
	return strings.ReplaceAll(o.JWT.Secret, o.JWT.Secret, "***")
}
