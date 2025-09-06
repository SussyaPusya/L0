package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/SussyaPusya/L0/pkg/test"
	"github.com/segmentio/kafka-go"
)

var messageCount = 10

func main() {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Topic:    "order",
		Balancer: &kafka.LeastBytes{},
	}

	for range messageCount {
		order := test.GenerateOrder()

		data, err := json.Marshal(order)
		if err != nil {

			continue
		}

		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Value: data,
				Key:   []byte(order.OrderUID),
			},
		)

		if err != nil {
			log.Println("failed to write messages", "error", err, "order", order)
		} else {
			log.Println("success", "order", order)
		}

		time.Sleep(500 * time.Millisecond)
	}
	log.Println("producer stopped")

}
