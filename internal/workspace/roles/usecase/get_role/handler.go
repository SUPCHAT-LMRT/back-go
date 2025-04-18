package get_role

import (
	"github.com/gin-gonic/gin"
)

type GetRoleHandler struct {
	useCase *GetRoleUseCase
}

func NewGetRoleHandler(useCase *GetRoleUseCase) *GetRoleHandler {
	return &GetRoleHandler{useCase: useCase}
}

func (h GetRoleHandler) Handle(c *gin.Context) {
	roleId := c.Param("role_id")
	if roleId == "" {
		c.JSON(400, gin.H{"error": "role_id is required"})
		return
	}

	role, err := h.useCase.Execute(c, roleId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"id":          role.Id,
		"name":        role.Name,
		"workspaceId": role.WorkspaceId,
		"permissions": role.Permissions,
	})
}
