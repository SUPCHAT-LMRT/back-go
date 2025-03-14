package search

import (
	"context"
	"fmt"
	meilisearch2 "github.com/meilisearch/meilisearch-go"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/utils"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
)

type SearchTermUseCaseDeps struct {
	uberdig.In
	Client                    *meilisearch.MeilisearchClient
	GetWorkspaceMemberUseCase *get_workpace_member.GetWorkspaceMemberUseCase
	GetUserByIdUseCase        *get_by_id.GetUserByIdUseCase
}

type SearchTermUseCase struct {
	deps SearchTermUseCaseDeps
}

func NewSearchTermUseCase(deps SearchTermUseCaseDeps) *SearchTermUseCase {
	return &SearchTermUseCase{deps: deps}
}

func (u SearchTermUseCase) Execute(ctx context.Context, term string, kind string) ([]*SearchResult, error) {
	var queries []*meilisearch2.SearchRequest
	// If no kind is specified, search in all indexes
	if kind == "message" {
		queries = []*meilisearch2.SearchRequest{u.messageQuery(term)}
	} else if kind == "channel" {
		queries = []*meilisearch2.SearchRequest{u.channelQuery(term)}
	} else {
		queries = []*meilisearch2.SearchRequest{u.messageQuery(term), u.channelQuery(term)}
	}

	searchResponse, err := u.deps.Client.Client.MultiSearchWithContext(ctx, &meilisearch2.MultiSearchRequest{
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
					// The message was found in a channel (in a workspace)
					result := message.SearchMessage{Data: message.SearchMessageChannelData{}}
					err = utils.Decode(hitMap["_formatted"].(map[string]any), &result)
					if err != nil {
						return nil, err
					}

					data := result.Data.(message.SearchMessageChannelData)
					workspaceMember, err := u.deps.GetWorkspaceMemberUseCase.Execute(ctx, data.WorkspaceId, result.AuthorId)
					if err != nil {
						return nil, err
					}

					if workspaceMember.Pseudo == "" {
						user, err := u.deps.GetUserByIdUseCase.Execute(ctx, result.AuthorId)
						if err != nil {
							return nil, err
						}
						workspaceMember.Pseudo = user.FullName()
					}

					results = append(results, &SearchResult{
						Kind: SearchResultKindMessage,
						Data: &SearchResultMessage{
							Id:         result.Id,
							Content:    result.Content,
							AuthorId:   string(result.AuthorId),
							AuthorName: workspaceMember.Pseudo,
							Href:       fmt.Sprintf("/workspaces/%s/channels/%s", workspaceMember.WorkspaceId, data.ChannelId),
						},
					})
				}
			}

			if response.IndexUID == "channels" {
				result := channel.SearchChannel{}
				err = utils.Decode(hitMap["_formatted"].(map[string]any), &result)
				if err != nil {
					return nil, err
				}

				results = append(results, &SearchResult{
					Kind: SearchResultKindChannel,
					Data: &SearchResultChannel{
						Id:    result.Id,
						Name:  result.Name,
						Topic: result.Topic,
						Href:  fmt.Sprintf("/workspaces/%s/channels/%s", result.WorkspaceId, result.Id),
					},
				})
			}

		}
	}

	return results, nil
}

func (u SearchTermUseCase) messageQuery(term string) *meilisearch2.SearchRequest {
	return &meilisearch2.SearchRequest{
		IndexUID:              "messages",
		AttributesToHighlight: []string{"Content"},
		Query:                 term,
	}
}

func (u SearchTermUseCase) channelQuery(term string) *meilisearch2.SearchRequest {
	return &meilisearch2.SearchRequest{
		IndexUID:              "channels",
		AttributesToHighlight: []string{"Name", "Topic"},
		Query:                 term,
	}
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
	Id    string `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
	Href  string `json:"href"`
}

type SearchResultMessage struct {
	Id         string `json:"id"`
	Content    string `json:"content"`
	AuthorId   string `json:"authorId"`
	AuthorName string `json:"authorName"`
	Href       string `json:"href"`
}
