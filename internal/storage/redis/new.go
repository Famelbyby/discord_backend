package redis

import (
	"errors"

	"github.com/caarlos0/env"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Password string `env:"REDIS_PASSWORD"`
}

func NewClient() (*redis.Client, error) {
	cfg := RedisConfig{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	emptyConfig := RedisConfig{}
	if cfg == emptyConfig {
		return nil, errors.New("[redis newClient] redis config is empty")
	}

	opt, err := redis.ParseURL("redis://default:" + cfg.Password + "@redis_container:6379/0")

	if err != nil {
		return nil, errors.New("[redis newClient] bad redis URL")
	}

	db := redis.NewClient(opt)

	return db, nil
}
