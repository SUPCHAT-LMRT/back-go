package search

import (
	"context"
	"fmt"
	meilisearch2 "github.com/meilisearch/meilisearch-go"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	"github.com/supchat-lmrt/back-go/internal/utils"
)

type SearchTermUseCase struct {
	client *meilisearch.MeilisearchClient
}

func NewSearchTermUseCase(client *meilisearch.MeilisearchClient) *SearchTermUseCase {
	return &SearchTermUseCase{client: client}
}

func (u SearchTermUseCase) Execute(ctx context.Context, term string, kind string) ([]*SearchResult, error) {
	var queries []*meilisearch2.SearchRequest
	// If no kind is specified, search in all indexes
	if kind == "" {
		queries = []*meilisearch2.SearchRequest{
			{
				IndexUID:              "messages",
				AttributesToHighlight: []string{"Content"},
				Query:                 term,
			},
		}
	} else if kind == "message" {
		queries = []*meilisearch2.SearchRequest{
			{
				IndexUID:              "messages",
				AttributesToHighlight: []string{"Content"},
				Query:                 term,
			},
		}
	}

	searchResponse, err := u.client.Client.MultiSearchWithContext(ctx, &meilisearch2.MultiSearchRequest{
		Queries: queries,
	})
	if err != nil {
		return nil, err
	}

	var results []*SearchResult
	for _, response := range searchResponse.Results {
		for _, hit := range response.Hits {
			hitMap := hit.(map[string]interface{})

			if response.IndexUID == "messages" {
				switch hitMap["Kind"].(string) {
				case string(SearchResultKindChannel):
					result := message.SearchMessage{Data: message.SearchMessageChannelData{}}
					err = utils.Decode(hitMap, &result)
					if err != nil {
						fmt.Println(err)
						return nil, err
					}
					// TODO voir pour implémenter le _formatted du hitMap
					results = append(results, &SearchResult{
						Kind: SearchResultKindMessage,
						Data: &SearchResultMessage{
							Id:         result.Id,
							Content:    result.Content,
							AuthorName: result.AuthorId.String(),
							Href:       fmt.Sprintf("/workspaces/%s/channels/%s", "monworkspaceid", result.Data.(message.SearchMessageChannelData).ChannelId),
						},
					})
				}
			}
		}
	}

	return results, nil
}

type SearchResult struct {
	Kind SearchResultKind `json:"kind"`
	Data any              `json:"data"`
}

type SearchResultKind string

const (
	SearchResultKindChannel SearchResultKind = "channel"
	SearchResultKindMessage SearchResultKind = "message"
)

type SearchResultChannel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SearchResultMessage struct {
	Id         string `json:"id"`
	Content    string `json:"content"`
	AuthorName string `json:"authorName"`
	Href       string `json:"href"`
}
