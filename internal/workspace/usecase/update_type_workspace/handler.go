package update_type_workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type UpdateTypeWorkspaceHandler struct {
	useCase *UpdateTypeWorkspaceUseCase
}

func NewUpdateTypeWorkspaceHandler(
	useCase *UpdateTypeWorkspaceUseCase,
) *UpdateTypeWorkspaceHandler {
	return &UpdateTypeWorkspaceHandler{useCase: useCase}
}

func (h *UpdateTypeWorkspaceHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(400, gin.H{"error": "workspace_id is required"})
		return
	}

	var request struct {
		Type string `json:"type"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid request"})
		return
	}

	err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId), entity.WorkspaceType(request.Type))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}
