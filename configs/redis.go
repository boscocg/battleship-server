package config

import (
	"context"
	"fmt"
	"log"
	"time"

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

func GetTimeLimit() time.Duration {
	exp := time.Hour // Default to 1 hour expiration
	// Get expiration from environment variable (minutes)
	if envTimeLimit := GetEnv("TIME_LIMIT"); envTimeLimit != "" {
		parsedTimeLimit, err := time.ParseDuration(envTimeLimit)
		if err != nil {
			log.Printf("Invalid TIME_LIMIT format, using default: %v", err)
			return time.Hour
		}
		exp = parsedTimeLimit
	}
	return exp
}
