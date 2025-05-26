package update_info_workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type UpdateInfoWorkspacesHandler struct {
	useCase *UpdateInfoWorkspacesUseCase
}

func NewUpdateInfoWorkspacesHandler(
	useCase *UpdateInfoWorkspacesUseCase,
) *UpdateInfoWorkspacesHandler {
	return &UpdateInfoWorkspacesHandler{useCase: useCase}
}

func (h *UpdateInfoWorkspacesHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	var request struct {
		Name  string `json:"name"`
		Topic string `json:"topic"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid request"})
		return
	}

	err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId), request.Name, request.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
