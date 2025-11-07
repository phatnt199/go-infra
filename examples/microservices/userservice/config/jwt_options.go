package config

import (
	"time"

	"github.com/iancoleman/strcase"
	"github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	"github.com/phatnt199/go-infra/pkg/crypto"
	"github.com/phatnt199/go-infra/pkg/reflection/typemapper"
)

type JWTOptions struct {
	Secret                 string `mapstructure:"secret" env:"JWT_SECRET" default:"your-secret-key-change-in-production"`
	AccessTokenExpiration  int64  `mapstructure:"accessTokenExpiration" env:"JWT_ACCESS_TOKEN_EXPIRATION" default:"3600"`
	RefreshTokenExpiration int64  `mapstructure:"refreshTokenExpiration" env:"JWT_REFRESH_TOKEN_EXPIRATION" default:"86400"`
}

func NewJWTOptions(environment environment.Environment) (*JWTOptions, error) {
	optionName := strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[JWTOptions]())

	cfg, err := config.BindConfigKey[*JWTOptions](optionName, environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (o *JWTOptions) ToJWTConfig() *crypto.JWTConfig {
	return &crypto.JWTConfig{
		Secret:             o.Secret,
		Algorithm:          crypto.AlgorithmHS256,
		AccessTokenExpiry:  time.Duration(o.AccessTokenExpiration) * time.Second,
		RefreshTokenExpiry: time.Duration(o.RefreshTokenExpiration) * time.Second,
	}
}
