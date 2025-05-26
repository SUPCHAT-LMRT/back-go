package reoder_channels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ReorderChannelRequest struct {
	Id       string `json:"id"`
	NewIndex int    `json:"newIndex"`
}

type ReorderChannelHandler struct {
	reorderUseCase *ReorderChannelsUseCase
}

func NewReorderChannelHandler(reorderUseCase *ReorderChannelsUseCase) *ReorderChannelHandler {
	return &ReorderChannelHandler{reorderUseCase: reorderUseCase}
}

func (h *ReorderChannelHandler) Handle(c *gin.Context) {
	var requests []ReorderChannelRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var inputs []ReorderChannelsInput
	for _, req := range requests {
		inputs = append(inputs, ReorderChannelsInput{
			ChannelId: entity.ChannelId(req.Id),
			NewIndex:  req.NewIndex,
		})
	}

	if err := h.reorderUseCase.ExecuteBulk(c, inputs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
