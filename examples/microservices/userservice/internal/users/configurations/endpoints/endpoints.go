package endpoints

import "github.com/phatnt199/go-infra/pkg/core/web/route"

func RegisterEndpoints(endpoints []route.Endpoint) error {
	for _, endpoint := range endpoints {
		endpoint.MapEndpoint()
	}
	return nil
}
