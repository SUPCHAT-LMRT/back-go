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

// Handle récupère les détails d'un rôle spécifique
// @Summary Détails d'un rôle
// @Description Récupère les informations détaillées d'un rôle spécifique dans un espace de travail
// @Tags workspace,role
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param role_id path string true "ID du rôle à récupérer"
// @Success 200 {object} map[string]string "Détails du rôle"
// @Failure 400 {object} map[string]string "ID du rôle manquant"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour voir les rôles"
// @Failure 404 {object} map[string]string "Rôle non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération du rôle"
// @Router /api/workspaces/{workspace_id}/roles/{role_id} [get]
// @Security ApiKeyAuth
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
