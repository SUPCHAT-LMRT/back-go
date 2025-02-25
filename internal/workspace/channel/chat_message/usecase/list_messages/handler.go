package list_messages

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
	"net/http"
	"time"
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

func NewListChannelMessagesHandler(deps ListChannelMessagesHandlerDeps) *ListChannelMessagesHandler {
	return &ListChannelMessagesHandler{deps: deps}
}

func (h *ListChannelMessagesHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspaceId")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "workspaceId is required"})
		return
	}

	channelId := c.Param("channelId")
	if channelId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "channelId is required"})
		return
	}

	channelMessages, err := h.deps.UseCase.Execute(c, channel_entity.ChannelId(channelId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := make([]ChannelMessageResponse, len(channelMessages))
	for i, message := range channelMessages {
		response[i] = ChannelMessageResponse{
			Id:        message.Id.String(),
			ChannelId: message.ChannelId.String(),
			Content:   message.Content,
			Author:    ChannelMessageAuthorResponse{},
			CreatedAt: message.CreatedAt,
		}

		user, err := h.deps.GetUserByIdUseCase.Execute(c, message.AuthorId)
		if err != nil {
			continue
		}

		member, err := h.deps.GetWorkspaceMemberUseCase.Execute(c, entity.WorkspaceId(workspaceId), message.AuthorId)
		if err != nil {
			continue
		}

		response[i].Author = ChannelMessageAuthorResponse{
			UserId:            user.Id.String(),
			WorkspaceMemberId: member.Id.String(),
			WorkspacePseudo:   member.Pseudo,
		}
	}

	c.JSON(http.StatusOK, response)
}

type ChannelMessageResponse struct {
	Id        string                       `json:"id"`
	ChannelId string                       `json:"channelId"`
	Content   string                       `json:"content"`
	Author    ChannelMessageAuthorResponse `json:"author"`
	CreatedAt time.Time                    `json:"createdAt"`
}

type ChannelMessageAuthorResponse struct {
	UserId            string `json:"userId"`
	WorkspaceMemberId string `json:"workspaceMemberId"`
	WorkspacePseudo   string `json:"workspacePseudo"`
}
