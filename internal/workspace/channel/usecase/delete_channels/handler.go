package delete_channels

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"net/http"
)

type DeleteChannelHandler struct {
	useCase *DeleteChannelUseCase
}

func NewDeleteChannelHandler(useCase *DeleteChannelUseCase) *DeleteChannelHandler {
	return &DeleteChannelHandler{useCase: useCase}
}

func (h *DeleteChannelHandler) Handle(c *gin.Context) {
	channelId := c.Param("channel_id")
	if channelId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel_id is required"})
		return
	}

	err := h.useCase.Execute(c, entity.ChannelId(channelId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
