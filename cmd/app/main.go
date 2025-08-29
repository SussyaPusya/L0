package main

import (
	"context"

	"github.com/SussyaPusya/L0/internal/config"
	"github.com/SussyaPusya/L0/internal/service"
	"github.com/SussyaPusya/L0/internal/transport/kafk"
)

func main() {

	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		//логи

	}

	svc := service.NewService()

	consum := kafk.NewConsumer(&cfg.Kafka, svc)

	consum.Consume(ctx)
}
