package user

import (
	"context"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	meilisearch2 "github.com/meilisearch/meilisearch-go"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type MeilisearchSearchUserSyncManager struct {
	createCache *lru.Cache[user_entity.UserId, *SearchUser]
	deleteCache *lru.Cache[user_entity.UserId, struct{}]
	client      *meilisearch.MeilisearchClient
	logger      logger.Logger
}

func NewMeilisearchSearchUserSyncManager(
	client *meilisearch.MeilisearchClient,
	logger logger.Logger,
) (SearchUserSyncManager, error) {
	createCache, err := lru.New[user_entity.UserId, *SearchUser](1000)
	if err != nil {
		return nil, err
	}

	deleteCache, err := lru.New[user_entity.UserId, struct{}](1000)
	if err != nil {
		return nil, err
	}

	return &MeilisearchSearchUserSyncManager{
		createCache: createCache,
		deleteCache: deleteCache,
		client:      client,
		logger:      logger,
	}, nil
}

func (m MeilisearchSearchUserSyncManager) CreateIndexIfNotExists(ctx context.Context) error {
	createdIndexTask, err := m.client.Client.CreateIndexWithContext(ctx, &meilisearch2.IndexConfig{
		Uid:        "users",
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
					"FirstName",
					"LastName",
					"Email",
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

func (m MeilisearchSearchUserSyncManager) AddUser(ctx context.Context, user *SearchUser) error {
	m.createCache.Add(user.Id, user)
	return nil
}

func (m MeilisearchSearchUserSyncManager) UpdateUser(ctx context.Context, user *SearchUser) error {
	m.createCache.Add(user.Id, user)
	return nil
}

func (m MeilisearchSearchUserSyncManager) RemoveUser(ctx context.Context, user *SearchUser) error {
	// Remove from main cache if exists
	m.createCache.Remove(user.Id)

	// Add to delete cache
	m.deleteCache.Add(user.Id, struct{}{})

	m.logger.Info().
		Str("user_id", user.Id.String()).
		Msg("User marked for deletion")

	return nil
}

func (m MeilisearchSearchUserSyncManager) Sync(ctx context.Context) {
	// Handle additions/updates
	var docs []*SearchUser
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
			Msg("Syncing users to Meilisearch")
		task, err := m.client.Client.Index("users").AddDocuments(docs)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to sync users to Meilisearch")
			return
		}

		// Wait for the task to complete
		taskInfo, err := m.client.Client.WaitForTask(task.TaskUID, 0)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to complete user sync task")
			return
		}
		if taskInfo.Status != meilisearch2.TaskStatusSucceeded {
			m.logger.Error().
				Str("status", string(taskInfo.Status)).
				Int("task_uid", int(taskInfo.TaskUID)).
				Msg("User sync task failed")
			return
		} else {
			m.createCache.Purge() // Clear only after successful sync
		}
	}

	// Sync deletions
	if len(deleteIds) > 0 {
		m.logger.Info().
			Int("count", len(deleteIds)).
			Msg("Syncing user deletions to Meilisearch")
		task, err := m.client.Client.Index("users").DeleteDocuments(deleteIds)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to sync user deletions to Meilisearch")
			return
		}

		// Wait for the deletion task to complete
		taskInfo, err := m.client.Client.WaitForTask(task.TaskUID, 0)
		if err != nil {
			m.logger.Error().
				Err(err).
				Msg("Failed to complete user deletion task")
			return
		}

		if taskInfo.Status != meilisearch2.TaskStatusSucceeded {
			m.logger.Error().
				Str("status", string(taskInfo.Status)).
				Int("task_uid", int(taskInfo.TaskUID)).
				Msg("User deletion task failed")
			return
		} else {
			m.deleteCache.Purge() // Clear only after successful sync
		}
	}
}

func (m MeilisearchSearchUserSyncManager) SyncLoop(ctx context.Context) {
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
