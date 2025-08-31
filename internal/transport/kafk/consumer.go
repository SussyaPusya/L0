package kafk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SussyaPusya/L0/internal/config"
	"github.com/SussyaPusya/L0/internal/dto"
	"github.com/segmentio/kafka-go"
)

type Service interface {
	CreateOrder(ctx context.Context, order *dto.Order) error
}
type Consumer struct {
	Reader  *kafka.Reader
	service Service
}

func NewConsumer(cfg *config.Kafka, s Service) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		GroupID:  "consumer-group-id",
		Topic:    "order",
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		Reader:  r,
		service: s,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {

	for {
		m, err := c.Reader.ReadMessage(context.Background())
		if err != nil {
			//логи нужны
			fmt.Println(err)
		}

		var order *dto.Order
		if err = json.Unmarshal(m.Value, &order); err != nil {
			fmt.Println(err)
			return err

		}

		if err = c.service.CreateOrder(ctx, order); err != nil {
			// логи
		}

		c.Reader.CommitMessages(ctx, m)

	}

}

func (c *Consumer) Shutdown() {

	c.Reader.Close()
}
