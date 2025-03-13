package main

import (
	"context"
	meilisearch2 "github.com/meilisearch/meilisearch-go"
	"github.com/supchat-lmrt/back-go/cmd/di"
	"github.com/supchat-lmrt/back-go/internal/gin"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/s3"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	"github.com/supchat-lmrt/back-go/internal/websocket"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	diContainer := di.NewDi()
	appContext := context.Background()

	var logg logger.Logger
	if err := diContainer.Invoke(func(logger logger.Logger) {
		logg = logger
	}); err != nil {
		log.Fatalf("Unable to invoke : %s", err.Error())
	}

	invokeFatal(logg, diContainer, func(httpServer gin.GinRouter) {
		httpServer.AddCorsHeaders()
		httpServer.RegisterRoutes()
	})

	logg.Info().Msg("Starting app...")

	invokeFatal(logg, diContainer, func(client *s3.S3Client) {
		logg.Info().Msg("Creating buckets...")
		bucketsToCreate := []string{"workspaces-icons", "workspaces-banners", "users-avatars", "messages-files"}

		bucketsCreated := make([]string, 0, len(bucketsToCreate))
		for _, bucket := range bucketsToCreate {
			created, err := client.CreateBucketIfNotExist(appContext, bucket)
			if err != nil {
				logg.Fatal().Err(err).Str("bucket", bucket).Msg("Unable to create bucket")
			}

			if created {
				bucketsCreated = append(bucketsCreated, bucket)
			}
		}

		logg.Info().Any("bucketsCreated", bucketsCreated).Msg("Buckets created!")
	})
	invokeFatal(logg, diContainer, func(client *mongo.Client) {
		// Create the time series collection "workspace_message_sent_ts" if it doesn't exist
		err := client.Client.Database("supchat").CreateCollection(appContext, "workspace_message_sent_ts", options.CreateCollection().
			SetTimeSeriesOptions(options.TimeSeries().
				SetTimeField("sent_at").
				SetMetaField("metadata").
				SetGranularity("minutes"),
			))
		if err != nil {
			if !strings.HasPrefix(err.Error(), "(NamespaceExists)") {
				logg.Fatal().Err(err).Msg("Unable to create collection")
			}
		}

		logg.Info().Str("collection", "workspace_message_sent_ts").Msg("Time-Series Collection created!")
	})

	invokeFatal(logg, diContainer, func(client *meilisearch.MeilisearchClient) {
		createdIndexTask, err := client.Client.CreateIndexWithContext(appContext, &meilisearch2.IndexConfig{
			Uid:        "messages",
			PrimaryKey: "Id",
		})
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to create index")
		}

		cancellableCtx, cancel := context.WithTimeout(appContext, 5*time.Second)
		defer cancel()

		task, err := client.Client.TaskReader().WaitForTaskWithContext(cancellableCtx, createdIndexTask.TaskUID, 0)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to wait for task")
		}

		if task.Status == meilisearch2.TaskStatusFailed {
			if task.Error.Code != "index_already_exists" {
				logg.Error().
					Str("status", string(task.Status)).
					Int("task_uid", int(task.TaskUID)).
					Str("details", task.Error.Code).
					Msg("Unable to create index")
			}

			return
		}

		if task.Status == meilisearch2.TaskStatusSucceeded {
			logg.Info().Str("uid", task.IndexUID).Msg("Index created!")
			updateSettingsTask, err := client.Client.Index(createdIndexTask.IndexUID).UpdateSettingsWithContext(appContext, &meilisearch2.Settings{
				DisplayedAttributes: []string{"*"},
				SearchableAttributes: []string{
					"Content",
				},
				FilterableAttributes: []string{
					"AuthorId",
					"Data.ChannelId",
					"Data.GroupId",
					"Data.OtherUserId",
					"CreatedAt",
					"UpdatedAt",
				},
				SortableAttributes: []string{
					"CreatedAt",
					"UpdatedAt",
				},
			})
			if err != nil {
				logg.Fatal().Err(err).Msg("Unable to update settings")
			}

			cancellableCtx, cancel = context.WithTimeout(appContext, 5*time.Second)
			defer cancel()
			task, err = client.Client.TaskReader().WaitForTaskWithContext(cancellableCtx, updateSettingsTask.TaskUID, 0)
			if err != nil {
				logg.Fatal().Err(err).Msg("Unable to wait for task")
			}

			if task.Status == meilisearch2.TaskStatusSucceeded {
				logg.Info().Msg("Settings updated!")
			} else {
				logg.Error().
					Str("status", string(updateSettingsTask.Status)).
					Int("task_uid", int(updateSettingsTask.TaskUID)).
					Msg("Unable to update settings")
			}
		}
	})

	invokeFatal(logg, diContainer, func(searchChannelMessageSyncManager message.SearchMessageSyncManager) {
		go searchChannelMessageSyncManager.SyncLoop(appContext)
	})

	go invokeFatal(logg, diContainer, runGinServer(logg))
	go invokeFatal(logg, diContainer, runWebsocketServer())

	logg.Info().Msg("App started!")

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, os.Interrupt, syscall.SIGTERM)

	<-signCh
	logg.Info().Msg("Shutting down app...")
}

func invokeFatal(logg logger.Logger, di *uberdig.Container, f any) {
	if err := di.Invoke(f); err != nil {
		logg.Fatal().Err(err).Msg("Unable to invoke")
	}
}

func runGinServer(logg logger.Logger) func(ginRouter gin.GinRouter) {
	return func(ginRouter gin.GinRouter) {
		err := ginRouter.Run()
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to run the router")
		}
	}
}

func runWebsocketServer() func(wsServer *websocket.WsServer) {
	return func(wsServer *websocket.WsServer) {
		wsServer.Run()
	}
}
