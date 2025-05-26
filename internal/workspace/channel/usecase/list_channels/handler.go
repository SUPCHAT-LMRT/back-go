package list_channels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type ListChannelsHandler struct {
	useCase *ListChannelsUseCase
}

func NewListChannelsHandler(useCase *ListChannelsUseCase) *ListChannelsHandler {
	return &ListChannelsHandler{useCase: useCase}
}

// TODO: filter out the channels that the user is not a member of
func (h *ListChannelsHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	channels, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]ListChannelResponse, len(channels))
	for i, channel := range channels {
		response[i] = ListChannelResponse{
			Id:    string(channel.Id),
			Name:  channel.Name,
			Topic: channel.Topic,
		}
	}

	c.JSON(http.StatusOK, response)
}

type ListChannelResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
}
