package consumer

import (
	"context"

	"github.com/phatnt199/go-infra/pkg/core/messaging/types"
)

type BusControl interface {
	// Start starts all consumers
	Start(ctx context.Context) error
	// Stop stops all consumers
	Stop() error

	IsConsumed(func(message types.IMessage))
}
