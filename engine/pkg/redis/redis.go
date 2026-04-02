package redis

import (
	"context"
	"fmt"

	"github.com/lieying/engine/internal/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func New(cfg config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil
}
