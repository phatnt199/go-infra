package infrastructure

import "github.com/phatnt199/go-infra/pkg/adapter/fxapp/contracts"

type InfrastructureConfigurator struct {
	contracts.Application
}

func NewInfrastructureConfigurator(
	fxapp contracts.Application,
) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: fxapp,
	}
}

func (ic *InfrastructureConfigurator) ConfigureInfrastructure() {
	ic.ResolveFunc(func() {

	})
}
