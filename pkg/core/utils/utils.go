package utils

import (
	"reflect"

	"github.com/phatnt199/go-infra/pkg/core/domain"
	"github.com/phatnt199/go-infra/pkg/core/events"
	typeMapper "github.com/phatnt199/go-infra/pkg/reflection/typemapper"

	"github.com/ahmetb/go-linq/v3"
)

func GetAllDomainEventTypes() []reflect.Type {
	var types []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).
		SelectManyT(func(i linq.KeyValue) linq.Query {
			return linq.From(i.Value)
		})
	d.ToSlice(&types)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[domain.IDomainEvent](types)
	linq.From(res).Distinct().ToSlice(&types)

	return types
}

func GetAllEventTypes() []reflect.Type {
	var types []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).
		SelectManyT(func(i linq.KeyValue) linq.Query {
			return linq.From(i.Value)
		})
	d.ToSlice(&types)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[events.IEvent](types)
	linq.From(res).Distinct().ToSlice(&types)

	return types
}
