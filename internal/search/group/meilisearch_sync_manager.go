package group

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	meilisearch2 "github.com/meilisearch/meilisearch-go"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
)

type (
	createLruCache = lru.Cache[group_entity.GroupId, *SearchGroup]
	deleteLruCache = lru.Cache[group_entity.GroupId, struct{}]
)

type MeilisearchSearchGroupSyncManager struct {
	createCache *createLruCache
	deleteCache *deleteLruCache
	client      *meilisearch.MeilisearchClient
	logger      logger.Logger
}

func NewMeilisearchSearchGroupSyncManager(
	client *meilisearch.MeilisearchClient,
	logg logger.Logger,
) (SearchGroupSyncManager, error) {
	createCache, err := lru.New[group_entity.GroupId, *SearchGroup](1000)
	if err != nil {
		return nil, err
	}

	deleteCache, err := lru.New[group_entity.GroupId, struct{}](1000)
	if err != nil {
		return nil, err
	}

	return &MeilisearchSearchGroupSyncManager{
		createCache: createCache,
		deleteCache: deleteCache,
		client:      client,
		logger:      logg,
	}, nil
}

//nolint:revive
func (m MeilisearchSearchGroupSyncManager) CreateIndexIfNotExists(ctx context.Context) error {
	createdIndexTask, err := m.client.Client.CreateIndexWithContext(ctx, &meilisearch2.IndexConfig{
		Uid:        "groups",
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
				},
				FilterableAttributes: []string{
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

func (m MeilisearchSearchGroupSyncManager) AddGroup(
	ctx context.Context,
	group *SearchGroup,
) error {
	m.createCache.Add(group.Id, group)
	return nil
}

func (m MeilisearchSearchGroupSyncManager) RemoveGroup(
	ctx context.Context,
	groupId group_entity.GroupId,
) error {
	// Remove from main cache if exists
	m.createCache.Remove(groupId)

	// Add to delete cache
	m.deleteCache.Add(groupId, struct{}{})

	m.logger.Info().
		Str("group_id", groupId.String()).
		Msg("Group marked for deletion")

	return nil
}

//nolint:revive
func (m MeilisearchSearchGroupSyncManager) Sync(ctx context.Context) {
	// Handle additions/updates
	var docs []*SearchGroup
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
			Msg("Syncing groups to Meilisearch")
		task, err := m.client.Client.Index("groups").AddDocuments(docs)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to sync groups to Meilisearch")
			return
		}

		// Wait for the task to complete
		taskInfo, err := m.client.Client.WaitForTask(task.TaskUID, 0)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to complete group sync task")
			return
		}
		if taskInfo.Status != meilisearch2.TaskStatusSucceeded {
			m.logger.Error().
				Str("status", string(taskInfo.Status)).
				Int("task_uid", int(taskInfo.TaskUID)).
				Msg("Group sync task failed")
			return
		}

		m.createCache.Purge() // Clear only after successful sync
	}

	// Sync deletions
	if len(deleteIds) > 0 {
		m.logger.Info().
			Int("count", len(deleteIds)).
			Msg("Syncing group deletions to Meilisearch")
		task, err := m.client.Client.Index("groups").DeleteDocuments(deleteIds)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to sync group deletions to Meilisearch")
			return
		}

		// Wait for the deletion task to complete
		taskInfo, err := m.client.Client.WaitForTask(task.TaskUID, 0)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to complete group deletion task")
			return
		}

		if taskInfo.Status != meilisearch2.TaskStatusSucceeded {
			m.logger.Error().
				Str("status", string(taskInfo.Status)).
				Int("task_uid", int(taskInfo.TaskUID)).
				Msg("Group deletion task failed")
			return
		}

		m.deleteCache.Purge() // Clear only after successful sync
	}
}

func (m MeilisearchSearchGroupSyncManager) SyncLoop(ctx context.Context) {
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
