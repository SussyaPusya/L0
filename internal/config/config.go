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
	Host     string `env:"PG_HOST"`
	Port     int    `env:"PG_PORT"`
	Database string `env:"PGPG_DATABASE_DATA"`
	User     string `env:"PG_USER"`
	Password string `env:"PG_PASS"`
	MaxConn  int32  `env:"PG_MAXCONN"`
	MinConn  int32  `env:"PG_MINCONN"`
}

type Redis struct {
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
