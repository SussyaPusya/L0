package main

import (
	"context"
	"fmt"

	"github.com/SussyaPusya/L0/internal/config"
	"github.com/SussyaPusya/L0/internal/repository"
	"github.com/SussyaPusya/L0/internal/service"
	"github.com/SussyaPusya/L0/internal/transport/kafk"
	"github.com/SussyaPusya/L0/internal/transport/rest"
	"github.com/SussyaPusya/L0/pkg/logger"
	"github.com/SussyaPusya/L0/pkg/postgres"
	"github.com/SussyaPusya/L0/pkg/redis"
	"go.uber.org/zap"
)

func main() {

	ctx := context.Background()
	loger, err := logger.NewLogger()
	if err != nil {
		panic(err)

	}

	cfg, err := config.NewConfig()
	if err != nil {
		loger.Error("failed to load config", zap.Error(err))
	}

	pg, err := postgres.NewPostgres(ctx, &cfg.Postgres)

	if err != nil {
		loger.Error("filed to connect to postgres", zap.Error(err))
	}

	cache, err := redis.NewRedis(&cfg.Redis, ctx)
	if err != nil {
		loger.Error("failed to connect to redis", zap.Error(err))
		fmt.Println(err)
	}

	repo := repository.NewRepository(pg, cache)

	svc := service.NewService(repo, loger)

	hanlrs := rest.NewHandlers(svc)
	router := rest.NewRouter(hanlrs)

	consumer := kafk.NewConsumer(&cfg.Kafka, svc)
	go router.Run()
	consumer.Consume(ctx)
}
