package meilisearch

import (
	"github.com/meilisearch/meilisearch-go"
	"os"
)

type MeilisearchClient struct {
	Client meilisearch.ServiceManager
}

func NewClient() (*MeilisearchClient, error) {
	client := meilisearch.New(os.Getenv("MEILISEARCH_HOST"), meilisearch.WithAPIKey(os.Getenv("MEILISEARCH_MASTER_KEY")))
	return &MeilisearchClient{Client: client}, nil
}
