package create_channel

import (
	"github.com/gin-gonic/gin"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"net/http"
)

type CreateChannelHandler struct {
	useCase *CreateChannelUseCase
}

func NewCreateChannelHandler(useCase *CreateChannelUseCase) *CreateChannelHandler {
	return &CreateChannelHandler{useCase: useCase}
}

func (h *CreateChannelHandler) Handle(c *gin.Context) {
	var req CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	err := h.useCase.Execute(c, &channel_entity.Channel{
		Name:        req.Name,
		Topic:       req.Topic,
		WorkspaceId: entity.WorkspaceId(workspaceId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

type CreateChannelRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Topic string `json:"topic"`
}
