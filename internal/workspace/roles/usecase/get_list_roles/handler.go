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

// Handle récupère la liste des rôles dans un espace de travail
// @Summary Liste des rôles
// @Description Récupère la liste de tous les rôles disponibles dans un espace de travail spécifié
// @Tags workspace,role
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Success 200 {array} map[string]string "Liste des rôles"
// @Failure 400 {object} map[string]string "ID de l'espace de travail manquant"
// @Failure 403 {object} map[string]string "Permissions insuffisantes ou utilisateur non membre de l'espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des rôles"
// @Router /api/workspaces/{workspace_id}/roles [get]
// @Security ApiKeyAuth
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
