package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/nfwGytautas/go-boilerplate/broker"
)

type ConsumerConfig struct {
	Group  string
	Topics []string
}

type consumerHandler[T any] struct {
	ctx context.Context
	c   broker.Consumer[T]
}

func NewConsumer[T any](ctx context.Context, broker *Broker, cfg ConsumerConfig, consumer broker.Consumer[T]) error {
	cg, err := sarama.NewConsumerGroupFromClient(cfg.Group, broker.client)
	if err != nil {
		return fmt.Errorf("create consumer group: %w", err)
	}
	broker.closables = append(broker.closables, cg)

	handler := consumerHandler[T]{
		ctx: ctx,
		c:   consumer,
	}

	for ctx.Err() == nil {
		if err := cg.Consume(ctx, cfg.Topics, &handler); err != nil {
			return err
		}
	}

	// The context is finished so no error
	return nil
}

func (h *consumerHandler[T]) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler[T]) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var data T
	for msg := range claim.Messages() {
		err := json.Unmarshal(msg.Value, &data)
		if err != nil {
			slog.Error(
				"broker_kafka_consume",
				"when", "unmarshal",
				"error", err,
			)
		}

		if err := h.c.Handle(h.ctx, data); err != nil {
			return err
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
