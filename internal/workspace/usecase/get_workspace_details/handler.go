package get_workspace_details

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetWorkspaceDetailsHandler struct {
	useCase *GetWorkspaceDetailsUseCase
}

func NewGetWorkspaceDetailsHandler(
	useCase *GetWorkspaceDetailsUseCase,
) *GetWorkspaceDetailsHandler {
	return &GetWorkspaceDetailsHandler{useCase: useCase}
}

func (h *GetWorkspaceDetailsHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	workspaceDetails, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, WorkspaceDetailsResponse{
		Id:            workspaceDetails.Id.String(),
		Name:          workspaceDetails.Name,
		Topic:         workspaceDetails.Topic,
		Type:          string(workspaceDetails.Type),
		MembersCount:  workspaceDetails.MembersCount,
		ChannelsCount: workspaceDetails.ChannelsCount,
		MessagesCount: workspaceDetails.MessagesCount,
	})
}

type WorkspaceDetailsResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Topic         string `json:"topic"`
	Type          string `json:"type"`
	MembersCount  uint   `json:"membersCount"`
	ChannelsCount uint   `json:"channelsCount"`
	MessagesCount uint   `json:"messagesCount"`
}
