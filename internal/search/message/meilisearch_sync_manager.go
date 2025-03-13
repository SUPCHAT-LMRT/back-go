package message

import (
	"context"
	lru "github.com/hashicorp/golang-lru/v2"
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

func (m MeilisearchSearchMessageSyncManager) AddMessage(ctx context.Context, message *SearchMessage) error {
	m.cache.Add(message.Id, message)
	return nil
}

func (m MeilisearchSearchMessageSyncManager) SyncLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
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
					continue
				}
				m.cache.Purge()
			}
		}
	}
}
