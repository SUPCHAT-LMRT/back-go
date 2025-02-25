package main

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/cmd/pubsubtest/messages"
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

	// Create a message
	var message messages.Message = messages.SendMessageToChannel{
		Sender: &messages.SendMessageToChannelSender{
			UserId:            "user-id",
			Pseudo:            "user-pseudo",
			WorkspaceMemberId: "workspace-member-id",
			WorkspacePseudo:   "workspace-pseudo",
		},
		ChannelId: "channel-id",
		Content:   "Hello, World!",
	}

	// Marshal the message
	marshal, err := json.Marshal(message)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to marshal message")
	}

	// Publish a message to the channel
	err = redisClient.Client.Publish(ctx, "ws-messages", marshal).Err()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to publish message")
	}
}
