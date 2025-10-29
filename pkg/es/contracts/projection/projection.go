package projection

import (
	"context"

	"github.com/phatnt199/go-infra/pkg/es/models"
)

type IProjection interface {
	ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error
}
