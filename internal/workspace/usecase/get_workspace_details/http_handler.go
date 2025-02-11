package get_workspace_details

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"net/http"
)

type GetWorkspaceDetailsHandler struct {
	useCase *GetWorkspaceDetailsUseCase
}

func NewGetWorkspaceDetailsHandler(useCase *GetWorkspaceDetailsUseCase) *GetWorkspaceDetailsHandler {
	return &GetWorkspaceDetailsHandler{useCase: useCase}
}

func (h *GetWorkspaceDetailsHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspaceId")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspaceId is required"})
		return
	}

	workspaceDetails, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, WorkspaceDetailsResponse{
		MembersCount:  workspaceDetails.MembersCount,
		ChannelsCount: workspaceDetails.ChannelsCount,
		MessagesCount: workspaceDetails.MessagesCount,
	})
}

type WorkspaceDetailsResponse struct {
	MembersCount  uint `json:"membersCount"`
	ChannelsCount uint `json:"channelsCount"`
	MessagesCount uint `json:"messagesCount"`
}
