package app

import (
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp/contracts"
)

type UserApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewUserApplicationBuilder() *UserApplicationBuilder {
	builder := &UserApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

func (b *UserApplicationBuilder) Build() *UserApplication {
	return NewUserApplication(
		b.GetProvides(),
		b.GetDecorates(),
		b.Options(),
		b.Logger(),
		b.Environment(),
	)
}
