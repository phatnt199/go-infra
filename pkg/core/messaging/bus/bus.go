package bus

import (
	consumer2 "github.com/phatnt199/go-infra/pkg/core/messaging/consumer"
	"github.com/phatnt199/go-infra/pkg/core/messaging/producer"
)

type Bus interface {
	producer.Producer
	consumer2.BusControl
	consumer2.ConsumerConnector
}
