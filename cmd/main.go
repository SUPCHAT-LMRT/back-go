package main

import (
	"context"
	"github.com/supchat-lmrt/back-go/cmd/di"
	"github.com/supchat-lmrt/back-go/internal/gin"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/s3"
	"github.com/supchat-lmrt/back-go/internal/websocket"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	diContainer := di.NewDi()

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
			created, err := client.CreateBucketIfNotExist(context.Background(), bucket)
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
		err := client.Client.Database("supchat").CreateCollection(context.Background(), "workspace_message_sent_ts", options.CreateCollection().
			SetTimeSeriesOptions(options.TimeSeries().
				SetTimeField("sent_at").
				SetMetaField("metadata").
				SetGranularity("minutes"),
			))
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to create collection")
		}

		logg.Info().Str("collection", "workspace_message_sent_ts").Msg("Time-Series Collection created!")
	})

	go invokeFatal(logg, diContainer, runGinServer(logg))
	go invokeFatal(logg, diContainer, runWebsocketServer(logg))

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

func runWebsocketServer(logg logger.Logger) func(wsServer *websocket.WsServer) {
	return func(wsServer *websocket.WsServer) {
		wsServer.Run()
	}
}
