package update_role

import (
	"github.com/gin-gonic/gin"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
)

type UpdateRoleHandler struct {
	useCase *UpdateRoleUseCase
}

func NewUpdateRoleHandler(useCase *UpdateRoleUseCase) *UpdateRoleHandler {
	return &UpdateRoleHandler{useCase: useCase}
}

func (h UpdateRoleHandler) Handle(c *gin.Context) {
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	roleId := c.Param("role_id")
	if roleId == "" {
		c.JSON(400, gin.H{"error": "role_id is required"})
		return
	}

	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(400, gin.H{"error": "workspace_id is required"})
		return
	}

	role := entity.Role{
		Id:          entity.RoleId(roleId),
		Name:        req.Name,
		WorkspaceId: workspace_entity.WorkspaceId(workspaceId),
		Permissions: req.Permissions,
	}

	err := h.useCase.Execute(c, role)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Permissions uint64 `json:"permissions" binding:"required"`
}
