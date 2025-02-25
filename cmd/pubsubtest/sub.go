package main

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger/zerolog"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"os"
)

func main() {
	os.Setenv("REDIS_URI", "redis://:password@localhost:6379")

	logger := zerolog.NewZerologLogger()
	redisClient, err := redis.NewClient()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to create redis client")
	}
	defer redisClient.Close()

	ctx := context.Background()
	err = redisClient.Client.Ping(ctx).Err()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to ping redis")
	}

	logger.Info().Msg("redis connected")

	// Subscribe
	pubsub := redisClient.Client.Subscribe(ctx, "ws-messages")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		logger.Info().Msgf("received message: %s", msg.Payload)
	}
}
