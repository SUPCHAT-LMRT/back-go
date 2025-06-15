package delete_role

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteRoleHandler struct {
	useCase *DeleteRoleUseCase
}

func NewDeleteRoleHandler(useCase *DeleteRoleUseCase) *DeleteRoleHandler {
	return &DeleteRoleHandler{useCase: useCase}
}

// Handle supprime un rôle existant dans un espace de travail
// @Summary Suppression d'un rôle
// @Description Supprime un rôle existant dans l'espace de travail spécifié
// @Tags workspace,role
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param role_id path string true "ID du rôle à supprimer"
// @Success 204 {string} string "Rôle supprimé avec succès"
// @Failure 400 {object} map[string]string "ID de rôle manquant"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour supprimer un rôle"
// @Failure 404 {object} map[string]string "Rôle non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la suppression du rôle"
// @Router /api/workspaces/{workspace_id}/roles/{role_id} [delete]
// @Security ApiKeyAuth
func (h DeleteRoleHandler) Handle(c *gin.Context) {
	roleId := c.Param("role_id")
	if roleId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_id is required"})
		return
	}

	err := h.useCase.Execute(c, roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
