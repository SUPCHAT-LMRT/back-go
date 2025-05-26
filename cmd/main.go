package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/supchat-lmrt/back-go/cmd/di"
	"github.com/supchat-lmrt/back-go/internal/gin"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/s3"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	"github.com/supchat-lmrt/back-go/internal/search/user"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/assign_job"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/crypt"
	"github.com/supchat-lmrt/back-go/internal/websocket"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
)

//nolint:revive
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
		bucketsToCreate := []string{
			"workspaces-icons",
			"workspaces-banners",
			"users-avatars",
			"messages-files",
		}

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
		err := client.Client.Database("supchat").
			CreateCollection(appContext, "workspace_message_sent_ts", options.CreateCollection().
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

		logg.Info().
			Str("collection", "workspace_message_sent_ts").
			Msg("Time-Series Collection created!")
	})

	// Ensure the Admin role exists and users
	invokeFatal(logg, diContainer, func(
		jobRepo repository.JobRepository,
		userRepository user_repository.UserRepository,
		assignJobUseCase *assign_job.AssignJobUseCase,
		cryptStrategy crypt.CryptStrategy,
	) {
		// Assurer l'existence du rôle Admin
		createdAdminRole, err := jobRepo.EnsureAdminJobExists(appContext)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to ensure Admin role exists")
		}
		logg.Info().Msg("Admin role ensured!")

		// Assurer l'existence du rôle Manager
		createdManagerRole, err := jobRepo.EnsureManagerJobExists(appContext)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to ensure Manager role exists")
		}
		logg.Info().Msg("Manager role ensured!")

		// Vérifier les utilisateurs existants
		users, err := userRepository.List(appContext)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to list users")
		}

		if len(users) > 0 {
			logg.Info().Msg("Users found! No need to create default users.")
			return
		}

		// Créer un utilisateur par défaut
		logg.Info().Msg("No users found! Creating default users...")
		hashedPassword, err := cryptStrategy.Hash(os.Getenv("INIT_USER_PASSWORD"))
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to hash password")
		}

		createdUser := &user_entity.User{
			FirstName: os.Getenv("INIT_USER_FIRST_NAME"),
			LastName:  os.Getenv("INIT_USER_LAST_NAME"),
			Email:     os.Getenv("INIT_USER_EMAIL"),
			Password:  hashedPassword,
		}
		err = userRepository.Create(appContext, createdUser)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to create default user")
		}

		logg.Info().
			Str("email", os.Getenv("INIT_USER_EMAIL")).
			Msg("Default user created!")

		// Assigner le rôle Admin à l'utilisateur par défaut
		err = assignJobUseCase.Execute(appContext, createdAdminRole.Id, createdUser.Id)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to assign Admin role to default user")
		}
		logg.Info().
			Str("email", os.Getenv("INIT_USER_EMAIL")).
			Str("role", createdAdminRole.Name).
			Msg("Default user assigned to Admin role!")

		// Assigner le rôle Manager à l'utilisateur par défaut (optionnel)
		err = assignJobUseCase.Execute(appContext, createdManagerRole.Id, createdUser.Id)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to assign Manager role to default user")
		}
		logg.Info().
			Str("email", os.Getenv("INIT_USER_EMAIL")).
			Str("role", createdManagerRole.Name).
			Msg("Default user assigned to Manager role!")
	})

	// Create the Meilisearch indexes if they don't exist
	invokeFatal(logg, diContainer, func(
		searchChannelMessageSyncManager message.SearchMessageSyncManager,
		searchChannelSyncManager channel.SearchChannelSyncManager,
		searchUserSyncManager user.SearchUserSyncManager,
	) {
		err := searchChannelMessageSyncManager.CreateIndexIfNotExists(appContext)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to create index")
		}

		err = searchChannelSyncManager.CreateIndexIfNotExists(appContext)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to create index")
		}

		err = searchUserSyncManager.CreateIndexIfNotExists(appContext)
		if err != nil {
			logg.Fatal().Err(err).Msg("Unable to create index")
		}
	})

	invokeFatal(logg, diContainer, func(
		searchChannelMessageSyncManager message.SearchMessageSyncManager,
		searchChannelSyncManager channel.SearchChannelSyncManager,
		searchUserSyncManager user.SearchUserSyncManager,
	) {
		go searchChannelMessageSyncManager.SyncLoop(appContext)
		go searchChannelSyncManager.SyncLoop(appContext)
		go searchUserSyncManager.SyncLoop(appContext)
	})

	go invokeFatal(logg, diContainer, runGinServer(logg))
	go invokeFatal(logg, diContainer, runWebsocketServer())

	logg.Info().Msg("App started!")

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, os.Interrupt, syscall.SIGTERM)

	<-signCh
	logg.Info().Msg("Shutting down app...")
}

func invokeFatal(logg logger.Logger, diContainer *uberdig.Container, f any) {
	if err := diContainer.Invoke(f); err != nil {
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
