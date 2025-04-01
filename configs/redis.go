package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type Config struct {
	RedisClient *redis.Client
	Ctx         context.Context
}

var AppConfig Config

func InitializeConfig(ctx context.Context) {
	AppConfig = Config{
		RedisClient: NewRedisClient(ctx),
		Ctx:         ctx,
	}
}

func NewRedisClient(ctx context.Context) *redis.Client {
	config := RedisConfig{
		Addr:     GetEnv("REDIS_ADDR"),
		Password: GetEnv("REDIS_PASSWORD"),
	}

	opt, _ := redis.ParseURL("rediss://default:" + config.Password + "@" + config.Addr)
	rdb := redis.NewClient(opt)

	env := GetEnv("ENV")
	if env == "local" {
		rdb = redis.NewClient(&redis.Options{
			Addr:             config.Addr,
			Password:         config.Password,
			DB:               config.DB,
			DisableIndentity: true,
		})
	}

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully")
	return rdb
}

func WriteToRedis(key string, value string) error {
	client := AppConfig.RedisClient
	if client == nil {
		return fmt.Errorf("redisClient is not initialized")
	}

	err := client.Set(AppConfig.Ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func ReadFromRedis(key string) (string, error) {
	client := AppConfig.RedisClient
	if client == nil {
		return "", fmt.Errorf("redisClient is not initialized")
	}

	val, err := client.Get(AppConfig.Ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}
