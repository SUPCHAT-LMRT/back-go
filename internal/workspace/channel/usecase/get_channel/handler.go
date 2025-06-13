package get_channel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type GetChannelHandler struct {
	useCase *GetChannelUseCase
}

func NewGetChannelHandler(useCase *GetChannelUseCase) *GetChannelHandler {
	return &GetChannelHandler{useCase: useCase}
}

func (h *GetChannelHandler) Handle(c *gin.Context) {
	channelId := c.Param("channel_id")
	channel, err := h.useCase.Execute(c.Request.Context(), entity.ChannelId(channelId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         channel.Id,
		"name":       channel.Name,
		"topic":      channel.Topic,
		"kind":       channel.Kind,
		"isPrivate":  channel.IsPrivate,
		"members":    channel.Members,
		"workspace":  channel.WorkspaceId,
		"created_at": channel.CreatedAt,
		"updated_at": channel.UpdatedAt,
		"index":      channel.Index,
	})
}

type ChannelResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Topic       string `json:"topic"`
	WorkspaceId string `json:"workspaceId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
