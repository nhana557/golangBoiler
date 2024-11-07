package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client{
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	urlRedis := App.Config.GetString("redis.url")
	port := App.Config.GetInt("redis.port")
	password := App.Config.GetString("redis.password")
	redisAddr := fmt.Sprintf("%s:%d", urlRedis, port)
	
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: password,
		DB: 0,
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)
	
	return client
}