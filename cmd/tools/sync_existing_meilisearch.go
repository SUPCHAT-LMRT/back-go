package main

import (
	"context"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/dig"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/logger/zerolog"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	"github.com/supchat-lmrt/back-go/internal/search/user"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	user_repository_mongo "github.com/supchat-lmrt/back-go/internal/user/repository/mongo"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	channel_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	uberdig "go.uber.org/dig"
	"log"
)

func main() {
	diContainer := uberdig.New()
	providers := []dig.Provider{
		dig.NewProvider(zerolog.NewZerologLogger),
		dig.NewProvider(mongo.NewClient),
		dig.NewProvider(redis.NewClient),
		dig.NewProvider(meilisearch.NewClient),

		dig.NewProvider(repository.NewChannelMessageMapper),
		dig.NewProvider(repository.NewMongoChannelMessageRepository),

		dig.NewProvider(channel_repository.NewMongoChannelMapper),
		dig.NewProvider(channel_repository.NewMongoChannelRepository),

		dig.NewProvider(message.NewMeilisearchSearchMessageSyncManager),
		dig.NewProvider(channel.NewMeilisearchSearchChannelSyncManager),
		dig.NewProvider(user.NewMeilisearchSearchUserSyncManager),

		dig.NewProvider(user_repository_mongo.NewMongoUserMapper),
		dig.NewProvider(user_repository_mongo.NewMongoUserRepository),
	}
	for _, provider := range providers {
		if err := diContainer.Provide(provider.Constructor, provider.ProvideOptions...); err != nil {
			log.Fatalf("Unable to provide %s : %s", provider, err.Error())
		}
	}

	appContext := context.Background()

	var logg logger.Logger
	if err := diContainer.Invoke(func(logger logger.Logger) {
		logg = logger
	}); err != nil {
		log.Fatalf("Unable to invoke : %s", err.Error())
	}

	// Synchronize the existing data (channels & messages)
	invokeFatal(logg, diContainer, func(
		channelMessageRepository repository.ChannelMessageRepository,
		channelRepository channel_repository.ChannelRepository,
		searchChannelMessageSyncManager message.SearchMessageSyncManager,
		searchChannelSyncManager channel.SearchChannelSyncManager,
	) {
		for _, workspaceId := range []workspace_entity.WorkspaceId{"67a9f911a58c8c9073991a7d"} {
			channels, err := channelRepository.List(appContext, workspaceId)
			if err != nil {
				logg.Fatal().Err(err).Msg("Unable to list channels")
			}

			for _, chann := range channels {
				messages, err := channelMessageRepository.ListByChannelId(appContext, chann.Id, repository.ListByChannelIdQueryParams{Limit: 50000})
				if err != nil {
					logg.Fatal().Err(err).Msg("Unable to list messages")
				}

				err = searchChannelSyncManager.AddChannel(appContext, &channel.SearchChannel{
					Id:          chann.Id,
					Name:        chann.Name,
					Topic:       chann.Topic,
					Kind:        mapChannelKindToSearchResultChannelKind(chann.Kind),
					WorkspaceId: workspaceId,
					CreatedAt:   chann.CreatedAt,
					UpdatedAt:   chann.UpdatedAt,
				})
				if err != nil {
					logg.Fatal().Err(err).Msg("Unable to add channel")
				}
				fmt.Println("Inserted channel", chann.Id)

				for _, channelMessage := range messages {
					err = searchChannelMessageSyncManager.AddMessage(appContext, &message.SearchMessage{
						Id:       channelMessage.Id.String(),
						Content:  channelMessage.Content,
						AuthorId: channelMessage.AuthorId,
						Kind:     message.SearchMessageKindChannelMessage,
						Data: message.SearchMessageChannelData{
							ChannelId:   channelMessage.ChannelId,
							WorkspaceId: workspaceId,
						},
						CreatedAt: channelMessage.CreatedAt,
						UpdatedAt: channelMessage.UpdatedAt,
					})
					if err != nil {
						logg.Fatal().Err(err).Msg("Unable to add message")
					}
				}

				fmt.Println("Inserted messages", len(messages))
			}

			fmt.Println("Inserted workspace", workspaceId)
		}
	})

	// Synchronize the existing data (users)
	invokeFatal(logg, diContainer, func(
		userRepository user_repository.UserRepository,
		searchUserSyncManager user.SearchUserSyncManager,
	) {
		users, err := userRepository.List(appContext)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to list users")
		}

		for _, u := range users {
			err = searchUserSyncManager.AddUser(appContext, &user.SearchUser{
				Id:        u.Id,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			})
			if err != nil {
				logg.Fatal().Err(err).Msg("Unable to add user")
			}
		}

		fmt.Println("Inserted users", len(users))
	})

	// Run synchronization
	invokeFatal(logg, diContainer, func(
		searchChannelMessageSyncManager message.SearchMessageSyncManager,
		searchChannelSyncManager channel.SearchChannelSyncManager,
		searchUserSyncManager user.SearchUserSyncManager,
	) {
		searchChannelSyncManager.Sync(appContext)
		searchChannelMessageSyncManager.Sync(appContext)
		searchUserSyncManager.Sync(appContext)
	})
}

func mapChannelKindToSearchResultChannelKind(kind channel_entity.ChannelKind) channel.SearchChannelKind {
	switch kind {
	case channel_entity.ChannelKindText:
		return channel.SearchChannelKindText
	case channel_entity.ChannelKindVoice:
		return channel.SearchChannelKindVoice
	default:
		return channel.SearchChannelKindUnknown
	}
}

func invokeFatal(logg logger.Logger, di *uberdig.Container, f any) {
	if err := di.Invoke(f); err != nil {
		logg.Fatal().Err(err).Msg("Unable to invoke")
	}
}
