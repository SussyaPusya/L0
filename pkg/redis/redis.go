package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SussyaPusya/L0/internal/config"
	"github.com/SussyaPusya/L0/internal/dto"
	"github.com/redis/go-redis/v9"
)

type RedisChache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(config *config.Redis, ctx context.Context) (*RedisChache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       0,
	})

	pong := client.Ping(ctx)
	if pong.Err() != nil {
		return nil, pong.Err()
	}
	defaultTTL := time.Hour * time.Duration(config.TTL)

	return &RedisChache{
		client: client,
		ttl:    defaultTTL,
	}, nil
}

func (c *RedisChache) Set(ctx context.Context, order *dto.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, order.OrderUID, data, c.ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisChache) Get(ctx context.Context, orderUID string) (*dto.Order, error) {

	orderJSON, err := c.client.Get(ctx, orderUID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("order not found in cache: %w", err)
		}
		return nil, fmt.Errorf("failed to get order from cache: %w", err)
	}

	var order dto.Order
	err = json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return &order, nil
}
