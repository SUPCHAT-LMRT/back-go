package get_list_roles

import (
	"github.com/gin-gonic/gin"
)

type GetListRolesHandler struct {
	useCase *GetListRolesUseCase
}

func NewGetListRolesHandler(useCase *GetListRolesUseCase) *GetListRolesHandler {
	return &GetListRolesHandler{useCase: useCase}
}

func (h GetListRolesHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(400, gin.H{"error": "workspace_id is required"})
		return
	}

	roles, err := h.useCase.Execute(c, workspaceId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var response []gin.H
	for _, role := range roles {
		response = append(response, gin.H{
			"id":          role.Id,
			"name":        role.Name,
			"workspaceId": role.WorkspaceId,
			"permissions": role.Permissions,
			"color":       role.Color,
		})
	}

	c.JSON(200, response)
}
