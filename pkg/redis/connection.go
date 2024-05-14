package redis

import (
	"context"
	"spf-playlist/pkg/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Client *redis.Client
}

func NewRedis(cfg config.GlobalEnv) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: "",
		DB:       0,
	})

	return &Client{
		Client: client,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.Client.Ping(ctx).Result()

	return err
}

func (c *Client) Close() error {
	return c.Client.Close()
}
