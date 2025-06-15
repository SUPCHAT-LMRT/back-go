package assign_role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
)

type AssignRoleToUserHandler struct {
	useCase *AssignRoleToUserUsecase
}

func NewAssignRoleToUserHandler(useCase *AssignRoleToUserUsecase) *AssignRoleToUserHandler {
	return &AssignRoleToUserHandler{useCase: useCase}
}

// Handle attribue un rôle à un membre de l'espace de travail
// @Summary Attribution d'un rôle
// @Description Attribue un rôle à un membre spécifique de l'espace de travail
// @Tags workspace,role
// @Accept json
// @Produce json
// @Param body body AssignRoleToUserRequest true "Données d'attribution du rôle"
// @Success 200 {object} map[string]string "Rôle attribué avec succès"
// @Failure 400 {object} map[string]string "Requête invalide ou données incomplètes"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour attribuer des rôles"
// @Failure 404 {object} map[string]string "Utilisateur ou rôle non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de l'attribution du rôle"
// @Router /api/workspaces/{workspace_id}/roles/assign [post]
// @Security ApiKeyAuth
func (h AssignRoleToUserHandler) Handle(c *gin.Context) {
	var req AssignRoleToUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.useCase.Execute(
		c,
		entity2.WorkspaceMemberId(req.UserId),
		entity.RoleId(req.RoleId),
		entity.WorkspaceId(req.WorkspaceId),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

type AssignRoleToUserRequest struct {
	RoleId      string `json:"role_id"      binding:"required,min=1,max=100"`
	UserId      string `json:"user_id"      binding:"required,min=1,max=100"`
	WorkspaceId string `json:"workspace_id" binding:"required,min=1,max=100"`
}
