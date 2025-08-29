package kafk

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/SussyaPusya/L0/internal/config"
	"github.com/SussyaPusya/L0/internal/dto"
	"github.com/segmentio/kafka-go"
)

type Service interface {
	CreateOrder(order *dto.Order) error
}
type Consumer struct {
	Reader  *kafka.Reader
	service Service
}

func NewConsumer(cfg *config.Kafka, s Service) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		GroupID:  "consumer-group-id",
		Topic:    cfg.Topic,
		MaxBytes: 10e6, // 10MB
	})

	err := createTopicIfNotExists(cfg, 3, 1)
	if err != nil {
		fmt.Println(err)
	}

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

		if err = c.service.CreateOrder(order); err != nil {
			// логи
		}

		c.Reader.CommitMessages(ctx, m)

	}

}

func (c *Consumer) Shutdown() {

	c.Reader.Close()
}

func createTopicIfNotExists(cfg *config.Kafka, numPartitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}

	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             cfg.Topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
}
