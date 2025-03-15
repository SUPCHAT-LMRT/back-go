package message

import (
	"context"
	lru "github.com/hashicorp/golang-lru/v2"
	meilisearch2 "github.com/meilisearch/meilisearch-go"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	"time"
)

type MeilisearchSearchMessageSyncManager struct {
	cache  *lru.Cache[string, *SearchMessage]
	client *meilisearch.MeilisearchClient
	logger logger.Logger
}

func NewMeilisearchSearchMessageSyncManager(client *meilisearch.MeilisearchClient, logger logger.Logger) (SearchMessageSyncManager, error) {
	cache, err := lru.New[string, *SearchMessage](1000)
	if err != nil {
		return nil, err
	}

	return &MeilisearchSearchMessageSyncManager{
		cache:  cache,
		client: client,
		logger: logger,
	}, nil
}

func (m MeilisearchSearchMessageSyncManager) CreateIndexIfNotExists(ctx context.Context) error {
	createdIndexTask, err := m.client.Client.CreateIndexWithContext(ctx, &meilisearch2.IndexConfig{
		Uid:        "messages",
		PrimaryKey: "Id",
	})
	if err != nil {
		return err
	}

	cancellableCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	task, err := m.client.Client.TaskReader().WaitForTaskWithContext(cancellableCtx, createdIndexTask.TaskUID, 0)
	if err != nil {
		return err
	}

	if task.Status == meilisearch2.TaskStatusFailed {
		if task.Error.Code != "index_already_exists" {
			m.logger.Error().
				Str("status", string(task.Status)).
				Int("task_uid", int(task.TaskUID)).
				Str("details", task.Error.Code).
				Msg("Unable to create index")
			return err
		}

		return nil
	}

	if task.Status == meilisearch2.TaskStatusSucceeded {
		m.logger.Info().Str("uid", task.IndexUID).Msg("Index created!")
		updateSettingsTask, err := m.client.Client.Index(createdIndexTask.IndexUID).UpdateSettingsWithContext(ctx, &meilisearch2.Settings{
			DisplayedAttributes: []string{"*"},
			SearchableAttributes: []string{
				"Id",
				"Content",
			},
			FilterableAttributes: []string{
				"AuthorId",
				"Kind",
				"Data.ChannelId",
				"Data.WorkspaceId",
				"Data.GroupId",
				"Data.User1",
				"Data.User2",
				"CreatedAt",
				"UpdatedAt",
			},
			SortableAttributes: []string{
				"CreatedAt",
				"UpdatedAt",
			},
			RankingRules: []string{
				"attribute",
				"words",
				"typo",
				"proximity",
				"sort",
				"exactness",
			},
		})
		if err != nil {
			return err
		}

		cancellableCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		task, err = m.client.Client.TaskReader().WaitForTaskWithContext(cancellableCtx, updateSettingsTask.TaskUID, 0)
		if err != nil {
			return err
		}

		if task.Status == meilisearch2.TaskStatusSucceeded {
			return nil
		} else {
			m.logger.Error().
				Str("status", string(updateSettingsTask.Status)).
				Int("task_uid", int(updateSettingsTask.TaskUID)).
				Msg("Unable to update settings")
			return err
		}
	}

	return nil
}

func (m MeilisearchSearchMessageSyncManager) AddMessage(ctx context.Context, message *SearchMessage) error {
	m.cache.Add(message.Id, message)
	return nil
}

func (m MeilisearchSearchMessageSyncManager) Sync(ctx context.Context) {
	var docs []*SearchMessage
	for _, key := range m.cache.Keys() {
		if doc, ok := m.cache.Get(key); ok {
			docs = append(docs, doc)
		}
	}
	if len(docs) > 0 {
		m.logger.Info().
			Int("count", len(docs)).
			Msg("Syncing messages to Meilisearch")
		_, err := m.client.Client.Index("messages").AddDocuments(docs)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to sync messages to Meilisearch")
			return
		}
		m.cache.Purge()
	}
}

func (m MeilisearchSearchMessageSyncManager) SyncLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.Sync(ctx)
		}
	}
}
