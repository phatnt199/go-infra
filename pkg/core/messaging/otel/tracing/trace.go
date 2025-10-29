package tracing

import (
	"github.com/phatnt199/go-infra/pkg/otel/tracing"

	"go.opentelemetry.io/otel/trace"
)

var MessagingTracer trace.Tracer

func init() {
	MessagingTracer = tracing.NewAppTracer(
		"github.com/phatnt199/go-infra/pkg/messaging",
	) // instrumentation name
}
