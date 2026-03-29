package kafka

import (
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
)

// BrokerConfig config for the kafka broker
type BrokerConfig struct {
	Brokers []string
}

// Broker kafka broker wrapper able to spawn Consumers and Producers
type Broker struct {
	client sarama.Client

	closables []closable
}

type closable interface {
	Close() error
}

// NewBroker starts a new broker, if it fails it returns an error
func NewBroker(cfg BrokerConfig) (*Broker, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewClient: %w", err)
	}

	return &Broker{
		client:    client,
		closables: make([]closable, 0),
	}, nil
}

// Close close the broker connection
func (b *Broker) Close() {
	for _, c := range b.closables {
		err := c.Close()
		if err != nil {
			slog.Error(
				"broker_kafka_close",
				"error", err,
			)
		}
	}

	err := b.client.Close()
	if err != nil {
		slog.Error(
			"broker_kafka_close",
			"error", err,
		)
	}
}
