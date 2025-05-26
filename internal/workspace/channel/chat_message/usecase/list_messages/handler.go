package list_messages

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
)

type ListChannelMessagesHandlerDeps struct {
	uberdig.In
	UseCase                   *ListChannelMessagesUseCase
	GetUserByIdUseCase        *get_by_id.GetUserByIdUseCase
	GetWorkspaceMemberUseCase *get_workpace_member.GetWorkspaceMemberUseCase
}

type ListChannelMessagesHandler struct {
	deps ListChannelMessagesHandlerDeps
}

func NewListChannelMessagesHandler(
	deps ListChannelMessagesHandlerDeps,
) *ListChannelMessagesHandler {
	return &ListChannelMessagesHandler{deps: deps}
}

type MessageQuery struct {
	Limit           int       `form:"limit,default=20,max=100"`
	Before          time.Time `form:"before"`
	After           time.Time `form:"after"`
	AroundMessageId string    `form:"aroundMessageId"`
}

func (h *ListChannelMessagesHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "workspace_id is required"})
		return
	}

	channelId := c.Param("channel_id")
	if channelId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "channel_id is required"})
		return
	}

	var query MessageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channelMessages, err := h.deps.UseCase.Execute(
		c,
		channel_entity.ChannelId(channelId),
		QueryParams{
			Limit:           query.Limit,
			Before:          query.Before,
			After:           query.After,
			AroundMessageId: channel_message_entity.ChannelMessageId(query.AroundMessageId),
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := make([]ChannelMessageResponse, len(channelMessages))
	for i, message := range channelMessages {
		reactions := make([]ChannelMessageReactionResponse, len(message.Reactions))
		for j, reaction := range message.Reactions {

			reactionUsers := make([]ChannelMessageReactionUserResponse, len(reaction.UserIds))
			for k, userId := range reaction.UserIds {
				memberReacted, err := h.deps.GetWorkspaceMemberUseCase.Execute(
					c,
					entity.WorkspaceId(workspaceId),
					userId,
				)
				if err != nil {
					continue
				}

        reactionUsers[k] = ChannelMessageReactionUserResponse{
          Id:   userId.String(),
          Name: userReacted.FullName(),
        }
			}

			reactions[j] = ChannelMessageReactionResponse{
				Id:       reaction.Id.String(),
				Users:    reactionUsers,
				Reaction: reaction.Reaction,
			}
		}

		response[i] = ChannelMessageResponse{
			Id:        message.Id.String(),
			ChannelId: message.ChannelId.String(),
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
			Reactions: reactions,
		}

		member, err := h.deps.GetWorkspaceMemberUseCase.Execute(
			c,
			entity.WorkspaceId(workspaceId),
			message.AuthorId,
		)
		if err != nil {
			continue
		}

		user, err := h.deps.GetUserByIdUseCase.Execute(c, message.AuthorId)
		if err != nil {
			continue
		}

		response[i].Author = ChannelMessageAuthorResponse{
			UserId:            message.AuthorId.String(),
			WorkspaceMemberId: member.Id.String(),
			WorkspacePseudo:   user.FullName(),
		}
	}

	c.JSON(http.StatusOK, response)
}

type ChannelMessageResponse struct {
	Id        string                           `json:"id"`
	ChannelId string                           `json:"channelId"`
	Content   string                           `json:"content"`
	Author    ChannelMessageAuthorResponse     `json:"author"`
	CreatedAt time.Time                        `json:"createdAt"`
	Reactions []ChannelMessageReactionResponse `json:"reactions"`
}

type ChannelMessageAuthorResponse struct {
	UserId            string `json:"userId"`
	WorkspaceMemberId string `json:"workspaceMemberId"`
	WorkspacePseudo   string `json:"workspacePseudo"`
}

type ChannelMessageReactionResponse struct {
	Id       string                               `json:"id"`
	Users    []ChannelMessageReactionUserResponse `json:"users"`
	Reaction string                               `json:"reaction"`
}

type ChannelMessageReactionUserResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
