package config

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	"github.com/phatnt199/go-infra/pkg/reflection/typemapper"
)

type AppOptions struct {
	ServiceName  string `mapstructure:"serviceName" env:"ServiceName"`
	DeliveryType string `mapstructure:"deliveryType" env:"DeliveryType" default:"http"`
}

func NewAppOptions(environment environment.Environment) (*AppOptions, error) {
	optionName := strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[AppOptions]())

	cfg, err := config.BindConfigKey[*AppOptions](optionName, environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (o *AppOptions) GetServiceName() string {
	return o.ServiceName
}

func (o *AppOptions) GetServiceNameUpper() string {
	return strings.ToUpper(o.ServiceName)
}
