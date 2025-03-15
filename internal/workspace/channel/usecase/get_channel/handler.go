package get_channel

import (
	"github.com/gin-gonic/gin"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"net/http"
)

type GetChannelHandler struct {
	useCase *GetChannelUseCase
}

func NewGetChannelHandler(useCase *GetChannelUseCase) *GetChannelHandler {
	return &GetChannelHandler{useCase: useCase}
}

func (h *GetChannelHandler) Handle(c *gin.Context) {
	channelId := c.Param("channel_id")

	channel, err := h.useCase.Execute(c, channel_entity.ChannelId(channelId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ChannelResponse{
		Id:          channel.Id.String(),
		Name:        channel.Name,
		Topic:       channel.Topic,
		WorkspaceId: channel.WorkspaceId.String(),
		CreatedAt:   channel.CreatedAt.String(),
		UpdatedAt:   channel.UpdatedAt.String(),
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
