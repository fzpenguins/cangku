package cache

import (
	"context"
	"time"
	"videoweb/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		//Username:    "127.0.0.1:6379",
		Password:    config.RedisPassword,
		DB:          config.RedisDB,
		ReadTimeout: 150 * time.Second,
		DialTimeout: 150 * time.Second,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	RedisClient = client
}
