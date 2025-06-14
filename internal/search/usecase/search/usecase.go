package search

import (
	"context"
	"fmt"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	group_search "github.com/supchat-lmrt/back-go/internal/search/group"

	meilisearch2 "github.com/meilisearch/meilisearch-go"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	user_search "github.com/supchat-lmrt/back-go/internal/search/user"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/utils"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
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

//nolint:revive
func (u SearchTermUseCase) Execute(
	ctx context.Context,
	term string,
	kind string,
	userInitiator *user_entity.User,
) ([]*SearchResult, error) {
	var queries []*meilisearch2.SearchRequest
	// If no kind is specified, search in all indexes
	switch kind {
	case "message":
		queries = []*meilisearch2.SearchRequest{u.messageQuery(term, userInitiator.Id)}
	case "channel":
		queries = []*meilisearch2.SearchRequest{u.channelQuery(term)}
	case "user":
		queries = []*meilisearch2.SearchRequest{u.userQuery(term)}
	case "group":
		queries = []*meilisearch2.SearchRequest{u.groupQuery(term)}
	default:
		queries = []*meilisearch2.SearchRequest{
			u.messageQuery(term, userInitiator.Id),
			u.channelQuery(term),
			u.userQuery(term),
			u.groupQuery(term),
		}
	}

	searchResponse, err := u.deps.Client.Client.MultiSearchWithContext(
		ctx,
		&meilisearch2.MultiSearchRequest{
			Queries: queries,
		},
	)
	if err != nil {
		return nil, err
	}

	var results []*SearchResult
	for _, response := range searchResponse.Results {
		for _, hit := range response.Hits {
			hitMap, ok := hit.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("unexpected hit type: %T", hit)
			}

			if response.IndexUID == "messages" {
				//nolint:revive
				switch hitMap["Kind"].(string) {
				case string(message.SearchMessageKindChannelMessage):
					// The message was found in a channel (in a workspace)
					result := message.SearchMessage{Data: message.SearchMessageChannelData{}}
					err = utils.Decode(hitMap["_formatted"].(map[string]any), &result)
					if err != nil {
						return nil, err
					}

					data, ok := result.Data.(message.SearchMessageChannelData)
					if !ok {
						return nil, fmt.Errorf("unexpected data type: %T", result.Data)
					}

					workspaceMember, err := u.deps.GetWorkspaceMemberUseCase.Execute(
						ctx,
						data.WorkspaceId,
						result.AuthorId,
					)
					if err != nil {
						return nil, err
					}

					user, err := u.deps.GetUserByIdUseCase.Execute(ctx, result.AuthorId)
					if err != nil {
						return nil, err
					}

					results = append(results, &SearchResult{
						Kind: SearchResultKindMessage,
						Data: &SearchResultMessage{
							Id:         result.Id,
							Content:    result.Content,
							AuthorId:   result.AuthorId,
							AuthorName: user.FullName(),
							Href: fmt.Sprintf(
								"/workspaces/%s/channels/%s?aroundMessageId=%s",
								workspaceMember.WorkspaceId,
								data.ChannelId,
								result.Id,
							),
						},
					})
				case string(message.SearchMessageKindDirectMessage):
					// The message was found in a channel (in a workspace)
					result := message.SearchMessage{Data: message.SearchMessageDirectData{}}
					err = utils.Decode(hitMap["_formatted"].(map[string]any), &result)
					if err != nil {
						return nil, err
					}

					data, ok := result.Data.(message.SearchMessageDirectData)
					if !ok {
						return nil, fmt.Errorf("unexpected data type: %T", result.Data)
					}

					user, err := u.deps.GetUserByIdUseCase.Execute(ctx, result.AuthorId)
					if err != nil {
						return nil, err
					}

					var otherUserId user_entity.UserId
					if userInitiator.Id == data.User1 {
						otherUserId = data.User2
					} else {
						otherUserId = data.User1
					}

					results = append(results, &SearchResult{
						Kind: SearchResultKindMessage,
						Data: &SearchResultMessage{
							Id:         result.Id,
							Content:    result.Content,
							AuthorId:   result.AuthorId,
							AuthorName: user.FullName(),
							Href: fmt.Sprintf(
								"/chat/direct/%s?aroundMessageId=%s",
								otherUserId,
								result.Id,
							),
						},
					})
				case string(message.SearchMessageGroupMessage):
					// The message was found in a channel (in a workspace)
					result := message.SearchMessage{Data: message.SearchMessageGroupData{}}
					err = utils.Decode(hitMap["_formatted"].(map[string]any), &result)
					if err != nil {
						return nil, err
					}

					data, ok := result.Data.(message.SearchMessageGroupData)
					if !ok {
						return nil, fmt.Errorf("unexpected data type: %T", result.Data)
					}

					user, err := u.deps.GetUserByIdUseCase.Execute(ctx, result.AuthorId)
					if err != nil {
						return nil, err
					}

					results = append(results, &SearchResult{
						Kind: SearchResultKindMessage,
						Data: &SearchResultMessage{
							Id:         result.Id,
							Content:    result.Content,
							AuthorId:   result.AuthorId,
							AuthorName: user.FullName(),
							Href: fmt.Sprintf(
								"/chat/group/%s?aroundMessageId=%s",
								data.GroupId,
								result.Id,
							),
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
						Kind:  u.mapSearchResultChannelKindToChannelKind(result.Kind),
						Href: fmt.Sprintf(
							"/workspaces/%s/channels/%s",
							result.WorkspaceId,
							result.Id,
						),
					},
				})
			}

			if response.IndexUID == "users" {
				result := user_search.SearchUser{}
				highlightedResult := user_search.SearchUser{}
				err = utils.Decode(hitMap, &result)
				if err != nil {
					return nil, err
				}
				err = utils.Decode(hitMap["_formatted"].(map[string]any), &highlightedResult)
				if err != nil {
					return nil, err
				}

				results = append(results, &SearchResult{
					Kind: SearchResultKindUser,
					Data: &SearchResultUser{
						Id:                   result.Id,
						HighlightedFirstName: highlightedResult.FirstName,
						HighlightedLastName:  highlightedResult.LastName,
						HighlightedEmail:     highlightedResult.Email,
						FirstName:            result.FirstName,
						LastName:             result.LastName,
						Email:                result.Email,
						Href:                 fmt.Sprintf("/chat/direct/%s", result.Id),
					},
				})
			}

			if response.IndexUID == "groups" {
				result := group_search.SearchGroup{}
				highlightedResult := group_search.SearchGroup{}
				err = utils.Decode(hitMap, &result)
				if err != nil {
					return nil, err
				}
				err = utils.Decode(hitMap["_formatted"].(map[string]any), &highlightedResult)
				if err != nil {
					return nil, err
				}

				results = append(results, &SearchResult{
					Kind: SearchResultKindGroup,
					Data: &SearchResultGroup{
						Id:              result.Id,
						HighlightedName: highlightedResult.Name,
						Name:            result.Name,
						Href: fmt.Sprintf(
							"/chat/group/%s",
							result.Id,
						),
					},
				})
			}

		}
	}

	return results, nil
}

// TODO Check if the user has access to this channel (or workspace) in case of channel messages
// TODO Check if the user is in the group in case of group messages
func (u SearchTermUseCase) messageQuery(
	term string,
	userInitiator user_entity.UserId,
) *meilisearch2.SearchRequest {
	return &meilisearch2.SearchRequest{
		IndexUID:              "messages",
		AttributesToHighlight: []string{"Content"},
		Query:                 term,
		Filter: fmt.Sprintf(
			"(Kind = 'direct' AND (Data.User1 = '%s' OR Data.User2 = '%s')) OR Kind != 'direct'",
			userInitiator,
			userInitiator,
		),
	}
}

// TODO Check if the user has access to this channel (or workspace)
func (u SearchTermUseCase) channelQuery(term string) *meilisearch2.SearchRequest {
	return &meilisearch2.SearchRequest{
		IndexUID:              "channels",
		AttributesToHighlight: []string{"Name", "Topic"},
		Query:                 term,
	}
}

func (u SearchTermUseCase) userQuery(term string) *meilisearch2.SearchRequest {
	return &meilisearch2.SearchRequest{
		IndexUID:              "users",
		AttributesToHighlight: []string{"FirstName", "LastName", "Email"},
		Query:                 term,
	}
}

func (u SearchTermUseCase) groupQuery(term string) *meilisearch2.SearchRequest {
	return &meilisearch2.SearchRequest{
		IndexUID:              "groups",
		AttributesToHighlight: []string{"Name"},
		Query:                 term,
	}
}

func (u SearchTermUseCase) mapSearchResultChannelKindToChannelKind(
	kind channel.SearchChannelKind,
) channel_entity.ChannelKind {
	switch kind {
	case channel.SearchChannelKindText:
		return channel_entity.ChannelKindText
	case channel.SearchChannelKindVoice:
		return channel_entity.ChannelKindVoice
	default:
		return channel_entity.ChannelKindUnknown
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
	SearchResultKindUser    SearchResultKind = "user"
	SearchResultKindGroup   SearchResultKind = "group"
)

type SearchResultChannel struct {
	Id    channel_entity.ChannelId   `json:"id"`
	Name  string                     `json:"name"`
	Topic string                     `json:"topic"`
	Kind  channel_entity.ChannelKind `json:"kind"`
	Href  string                     `json:"href"`
}

type SearchResultMessage struct {
	Id         string             `json:"id"`
	Content    string             `json:"content"`
	AuthorId   user_entity.UserId `json:"authorId"`
	AuthorName string             `json:"authorName"`
	Href       string             `json:"href"`
}

type SearchResultUser struct {
	Id                   user_entity.UserId `json:"id"`
	HighlightedFirstName string             `json:"highlightedFirstName"`
	HighlightedLastName  string             `json:"highlightedLastName"`
	HighlightedEmail     string             `json:"highlightedEmail"`
	FirstName            string             `json:"firstName"`
	LastName             string             `json:"lastName"`
	Email                string             `json:"email"`
	Href                 string             `json:"href"`
}

type SearchResultGroup struct {
	Id              group_entity.GroupId `json:"id"`
	HighlightedName string               `json:"highlightedName"`
	Name            string               `json:"name"`
	Href            string               `json:"href"`
}
