package get_workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetWorkspaceHandler struct {
	useCase *GetWorkspaceUseCase
}

func NewGetWorkspaceHandler(useCase *GetWorkspaceUseCase) *GetWorkspaceHandler {
	return &GetWorkspaceHandler{useCase: useCase}
}

func (h *GetWorkspaceHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	workspace, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    workspace.Id,
		"name":  workspace.Name,
		"topic": workspace.Topic,
		"type":  workspace.Type,
	})
}
