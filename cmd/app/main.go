package main

import (
	"context"
	"fmt"

	"github.com/SussyaPusya/L0/internal/config"
	"github.com/SussyaPusya/L0/internal/repository"
	"github.com/SussyaPusya/L0/internal/service"
	"github.com/SussyaPusya/L0/internal/transport/kafk"
	"github.com/SussyaPusya/L0/internal/transport/rest"
	"github.com/SussyaPusya/L0/pkg/postgres"
)

func main() {

	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		//логи

	}

	pg, err := postgres.NewPostgres(ctx, &cfg.Postgres)

	if err != nil {
		fmt.Println(err)
	}

	repo := repository.NewRepository(pg)

	svc := service.NewService(repo)

	hanlrs := rest.NewHandlers(svc)
	router := rest.NewRouter(hanlrs)

	consumer := kafk.NewConsumer(&cfg.Kafka, svc)
	go router.Run()
	consumer.Consume(ctx)
}
