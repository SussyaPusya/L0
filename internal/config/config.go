package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres Postgres
	Redis    Redis
	Kafka    Kafka
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	Database string `env:"POSTGRES_DB"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASS"`
	MaxConn  int32  `env:"POSTGRES_MAXCONN"`
	MinConn  int32  `env:"POSTGRES_MINCONN"`
}

type Redis struct {
	TTL      int    `env:"REDIS_TTL"`
	Host     string `env:"REDIS_HOST"`
	Password string `env:"REDIS_PASS"`
	Port     int    `env:"REDIS_PORT"`
}

type Kafka struct {
	Host  string `env:"KAFKA_HOST"`
	Port  int    `env:"KAFKA_PORT"`
	Topic string `env:"KAFKA_TOPIC"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Println(err)
	}

	fmt.Println(cfg)

	return &cfg, nil
}
