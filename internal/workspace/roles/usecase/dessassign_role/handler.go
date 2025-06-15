package dessassign_role

import "github.com/gin-gonic/gin"

type DessassignRoleFromUserHandler struct {
	useCase *DessassignRoleFromUserUsecase
}

func NewDessassignRoleFromUserHandler(
	useCase *DessassignRoleFromUserUsecase,
) *DessassignRoleFromUserHandler {
	return &DessassignRoleFromUserHandler{useCase: useCase}
}

// Handle retire un rôle à un membre de l'espace de travail
// @Summary Retrait d'un rôle
// @Description Retire un rôle spécifique à un membre de l'espace de travail
// @Tags workspace,role
// @Accept json
// @Produce json
// @Param body body DessassignRoleFromUserRequest true "Données de retrait du rôle"
// @Success 200 {object} map[string]string "Rôle retiré avec succès"
// @Failure 400 {object} map[string]string "Requête invalide ou données incomplètes"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour gérer les rôles"
// @Failure 404 {object} map[string]string "Utilisateur ou rôle non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors du retrait du rôle"
// @Router /api/workspaces/{workspace_id}/roles/dessassign [post]
// @Security ApiKeyAuth
func (h DessassignRoleFromUserHandler) Handle(c *gin.Context) {
	var req DessassignRoleFromUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.useCase.Execute(c, req.UserId, req.RoleId, req.WorkspaceId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Role unassigned successfully"})
}

type DessassignRoleFromUserRequest struct {
	RoleId      string `json:"role_id"      binding:"required,min=1,max=100"`
	UserId      string `json:"user_id"      binding:"required,min=1,max=100"`
	WorkspaceId string `json:"workspace_id" binding:"required,min=1,max=100"`
}
