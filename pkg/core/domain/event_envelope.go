package domain

import (
	"github.com/phatnt199/go-infra/pkg/core/metadata"
)

type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}
