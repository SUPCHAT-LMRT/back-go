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

// Handle crée un nouveau rôle dans un espace de travail
// @Summary Création d'un rôle
// @Description Crée un nouveau rôle avec un nom et une couleur dans l'espace de travail spécifié
// @Tags workspace,role
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param body body CreateRoleRequest true "Données du rôle à créer"
// @Success 200 {object} map[string]string "Rôle créé avec succès"
// @Failure 400 {object} map[string]string "Requête invalide ou données manquantes"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour créer un rôle"
// @Failure 500 {object} map[string]string "Erreur lors de la création du rôle"
// @Router /api/workspaces/{workspace_id}/roles [post]
// @Security ApiKeyAuth
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
