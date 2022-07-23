package publisher

import (
	"strings"

	"github.com/1005281342/postgresql-cdc/internal/config"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

// CreatePublisher is a helper function that creates a Publisher
func CreatePublisher(eq config.EventQueen) message.Publisher {
	if eq.Type != "kafka" {
		panic("暂不支持:" + eq.Type)
	}

	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   strings.Split(eq.Hosts, ","),
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NopLogger{},
	)
	if err != nil {
		panic(err)
	}

	return kafkaPublisher
}
