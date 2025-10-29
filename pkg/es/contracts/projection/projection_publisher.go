package projection

import (
	"context"

	"github.com/phatnt199/go-infra/pkg/es/models"
)

type IProjectionPublisher interface {
	Publish(ctx context.Context, streamEvent *models.StreamEvent) error
}
