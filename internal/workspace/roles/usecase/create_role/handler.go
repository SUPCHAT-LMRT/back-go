package create_role

import (
	"github.com/gin-gonic/gin"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
)

type CreateRoleHandler struct {
	useCase *CreateRoleUseCase
}

func NewCreateRoleHandler(useCase *CreateRoleUseCase) *CreateRoleHandler {
	return &CreateRoleHandler{useCase: useCase}
}

func (h CreateRoleHandler) Handle(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(400, gin.H{"error": "workspace_id is required"})
		return
	}

	if req.Color == "" {
		req.Color = "#6366f1"
	}

	role := entity.Role{
		Name:        req.Name,
		WorkspaceId: workspace_entity.WorkspaceId(workspaceId),
		Color:       req.Color,
	}

	err := h.useCase.Execute(c, &role)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"id":          role.Id,
		"name":        role.Name,
		"workspaceId": role.WorkspaceId,
		"permissions": role.Permissions,
		"color":       role.Color,
	}
	c.JSON(200, response)
}

type CreateRoleRequest struct {
	Name  string `json:"name"  binding:"required,min=1,max=100"`
	Color string `json:"color" binding:"required,min=1,max=100"`
}
