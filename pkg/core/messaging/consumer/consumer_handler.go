package consumer

import (
	"context"

	"github.com/phatnt199/go-infra/pkg/core/messaging/types"
)

type ConsumerHandler interface {
	Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error
}
