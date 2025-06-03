package generate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace"
)

type CreateInviteLinkHandler struct {
	useCase             *InviteLinkUseCase
	getWorkspaceUseCase *get_workspace.GetWorkspaceUseCase
}

func NewCreateInviteLinkHandler(
	usecase *InviteLinkUseCase,
	getWorkspaceUseCase *get_workspace.GetWorkspaceUseCase,
) *CreateInviteLinkHandler {
	return &CreateInviteLinkHandler{useCase: usecase, getWorkspaceUseCase: getWorkspaceUseCase}
}

type CreateInviteLinkRequest struct {
	WorkspaceId entity.WorkspaceId `json:"workspaceId"`
}

func (h *CreateInviteLinkHandler) Handle(c *gin.Context) {
	var request CreateInviteLinkRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if request.WorkspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "WorkspaceId cannot be empty"})
		return
	}

	_, err := h.getWorkspaceUseCase.Execute(c, request.WorkspaceId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	inviteLink, err := h.useCase.Execute(c, request.WorkspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, inviteLink)
}
