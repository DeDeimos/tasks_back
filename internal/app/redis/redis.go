package redis

import (
	"awesomeProject/internal/app/config"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const servicePrefix = "awesome_service." // наш префикс сервиса

type Client struct {
	cfg    config.RedisConfig
	client *redis.Client
}

func New(cfg config.RedisConfig) (*Client, error) {
	client := &Client{}
	log.Println("start done")
	cfg.DialTimeout = 30 * time.Second
	cfg.ReadTimeout = 30 * time.Second

	client.cfg = cfg

	log.Println(client.cfg)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: "",
		DB:       0,
	})
	log.Println(redisClient)

	client.client = redisClient
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("can't ping redis: %w", err)
	}
	log.Println("client done")

	return client, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
