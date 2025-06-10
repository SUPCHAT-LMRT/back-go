package channel

import (
	"context"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	meilisearch2 "github.com/meilisearch/meilisearch-go"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type (
	createLruCache = lru.Cache[channel_entity.ChannelId, *SearchChannel]
	deleteLruCache = lru.Cache[channel_entity.ChannelId, struct{}]
)

type MeilisearchSearchChannelSyncManager struct {
	createCache *createLruCache
	deleteCache *deleteLruCache
	client      *meilisearch.MeilisearchClient
	logger      logger.Logger
}

func NewMeilisearchSearchChannelSyncManager(
	client *meilisearch.MeilisearchClient,
	logg logger.Logger,
) (SearchChannelSyncManager, error) {
	createCache, err := lru.New[channel_entity.ChannelId, *SearchChannel](1000)
	if err != nil {
		return nil, err
	}

	deleteCache, err := lru.New[channel_entity.ChannelId, struct{}](1000)
	if err != nil {
		return nil, err
	}

	return &MeilisearchSearchChannelSyncManager{
		createCache: createCache,
		deleteCache: deleteCache,
		client:      client,
		logger:      logg,
	}, nil
}

//nolint:revive
func (m MeilisearchSearchChannelSyncManager) CreateIndexIfNotExists(ctx context.Context) error {
	createdIndexTask, err := m.client.Client.CreateIndexWithContext(ctx, &meilisearch2.IndexConfig{
		Uid:        "channels",
		PrimaryKey: "Id",
	})
	if err != nil {
		return err
	}

	cancellableCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	task, err := m.client.Client.TaskReader().
		WaitForTaskWithContext(cancellableCtx, createdIndexTask.TaskUID, 0)
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
		updateSettingsTask, err := m.client.Client.Index(createdIndexTask.IndexUID).
			UpdateSettingsWithContext(ctx, &meilisearch2.Settings{
				DisplayedAttributes: []string{"*"},
				SearchableAttributes: []string{
					"Id",
					"Name",
					"Topic",
				},
				FilterableAttributes: []string{
					"Kind",
					"Topic",
					"Data.ChannelId",
					"Data.WorkspaceId",
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
		task, err = m.client.Client.TaskReader().
			WaitForTaskWithContext(cancellableCtx, updateSettingsTask.TaskUID, 0)
		if err != nil {
			return err
		}

		if task.Status == meilisearch2.TaskStatusSucceeded {
			return nil
		}

		m.logger.Error().
			Str("status", string(updateSettingsTask.Status)).
			Int("task_uid", int(updateSettingsTask.TaskUID)).
			Msg("Unable to update settings")
		return err
	}

	return nil
}

func (m MeilisearchSearchChannelSyncManager) AddChannel(
	ctx context.Context,
	channel *SearchChannel,
) error {
	m.createCache.Add(channel.Id, channel)
	return nil
}

func (m MeilisearchSearchChannelSyncManager) RemoveChannel(
	ctx context.Context,
	channelId channel_entity.ChannelId,
) error {
	// Remove from main cache if exists
	m.createCache.Remove(channelId)

	// Add to delete cache
	m.deleteCache.Add(channelId, struct{}{})

	m.logger.Info().
		Str("channel_id", channelId.String()).
		Msg("Channel marked for deletion")

	return nil
}

//nolint:revive
func (m MeilisearchSearchChannelSyncManager) Sync(ctx context.Context) {
	// Handle additions/updates
	var docs []*SearchChannel
	for _, key := range m.createCache.Keys() {
		if doc, ok := m.createCache.Get(key); ok {
			docs = append(docs, doc)
		}
	}

	// Handle deletions
	var deleteIds []string
	for _, key := range m.deleteCache.Keys() {
		deleteIds = append(deleteIds, key.String())
	}

	// Sync additions/updates
	if len(docs) > 0 {
		m.logger.Info().
			Int("count", len(docs)).
			Msg("Syncing channels to Meilisearch")
		task, err := m.client.Client.Index("channels").AddDocuments(docs)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to sync channels to Meilisearch")
			return
		}

		// Wait for the task to complete
		taskInfo, err := m.client.Client.WaitForTask(task.TaskUID, 0)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to complete channel sync task")
			return
		}
		if taskInfo.Status != meilisearch2.TaskStatusSucceeded {
			m.logger.Error().
				Str("status", string(taskInfo.Status)).
				Int("task_uid", int(taskInfo.TaskUID)).
				Msg("Channel sync task failed")
			return
		}

		m.createCache.Purge() // Clear only after successful sync
	}

	// Sync deletions
	if len(deleteIds) > 0 {
		m.logger.Info().
			Int("count", len(deleteIds)).
			Msg("Syncing channel deletions to Meilisearch")
		task, err := m.client.Client.Index("channels").DeleteDocuments(deleteIds)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to sync channel deletions to Meilisearch")
			return
		}

		// Wait for the deletion task to complete
		taskInfo, err := m.client.Client.WaitForTask(task.TaskUID, 0)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to complete channel deletion task")
			return
		}

		if taskInfo.Status != meilisearch2.TaskStatusSucceeded {
			m.logger.Error().
				Str("status", string(taskInfo.Status)).
				Int("task_uid", int(taskInfo.TaskUID)).
				Msg("Channel deletion task failed")
			return
		}

		m.deleteCache.Purge() // Clear only after successful sync
	}
}

func (m MeilisearchSearchChannelSyncManager) SyncLoop(ctx context.Context) {
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
