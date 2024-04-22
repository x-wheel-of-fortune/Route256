package myRedis

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	cs, ok := os.LookupEnv("REDIS_CONNECTION_STRING")
	if !ok {
		log.Println("No redis connection string")
	}
	opt, err := redis.ParseURL(cs)
	if err != nil {
		log.Println("Не удалось распознать строку подключения к Redis")
		return nil
	}
	return &Redis{
		redis.NewClient(opt),
	}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, time.Minute*10).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	res := r.client.Get(ctx, key)
	return res.Result()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
