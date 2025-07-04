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

// Handle met à jour un rôle existant dans l'espace de travail
// @Summary Mise à jour d'un rôle
// @Description Met à jour les informations d'un rôle existant (nom, couleur, permissions)
// @Tags workspace,role
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param role_id path string true "ID du rôle à mettre à jour"
// @Param body body UpdateRoleRequest true "Nouvelles informations du rôle"
// @Success 200 {string} string "Rôle mis à jour avec succès"
// @Failure 400 {object} map[string]string "Requête invalide ou données incomplètes"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour gérer les rôles"
// @Failure 404 {object} map[string]string "Rôle non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la mise à jour du rôle"
// @Router /api/workspaces/{workspace_id}/roles/{role_id} [put]
// @Security ApiKeyAuth
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
		Color:       req.Color,
	}

	err := h.useCase.Execute(c, role)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

type UpdateRoleRequest struct {
	Name        string `json:"name"        binding:"required,min=1,max=100"`
	Permissions uint64 `json:"permissions" binding:"min=0"`
	Color       string `json:"color"       binding:"hexcolor"`
}
