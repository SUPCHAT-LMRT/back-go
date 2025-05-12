package list_private_channels

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"net/http"
)

type GetPrivateChannelsHandler struct {
	useCase *GetPrivateChannelsUseCase
}

func NewGetPrivateChannelsHandler(useCase *GetPrivateChannelsUseCase) *GetPrivateChannelsHandler {
	return &GetPrivateChannelsHandler{useCase: useCase}
}

func (h *GetPrivateChannelsHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	userId := c.Param("user_id")

	if workspaceId == "" || userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id and user_id are required"})
		return
	}

	channels, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]GetPrivateChannelResponse, len(channels))
	for i, channel := range channels {
		response[i] = GetPrivateChannelResponse{
			Id:    string(channel.Id),
			Name:  channel.Name,
			Topic: channel.Topic,
		}
	}

	c.JSON(http.StatusOK, response)
}

type GetPrivateChannelResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
}
