package get_roles_for_member

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

type GetRolesForMemberHandler struct {
	useCase *GetRolesForMemberUsecase
}

func NewGetRolesForMemberHandler(useCase *GetRolesForMemberUsecase) *GetRolesForMemberHandler {
	return &GetRolesForMemberHandler{useCase: useCase}
}

// Handle récupère les rôles attribués à un membre spécifique
// @Summary Rôles d'un membre
// @Description Récupère tous les rôles attribués à un membre spécifique dans l'espace de travail
// @Tags workspace,role,member
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param user_id path string true "ID de l'utilisateur"
// @Success 200 {object} map[string]string "Liste des rôles du membre"
// @Failure 400 {object} map[string]string "ID de l'espace de travail ou de l'utilisateur manquant"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour gérer les rôles"
// @Failure 404 {object} map[string]string "Membre non trouvé dans l'espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des rôles"
// @Router /api/workspaces/{workspace_id}/roles/members/{user_id} [get]
// @Security ApiKeyAuth
func (h GetRolesForMemberHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	userId := c.Param("user_id")

	roles, err := h.useCase.Execute(
		c,
		entity.WorkspaceId(workspaceId),
		entity2.WorkspaceMemberId(userId),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"roles": roles})
}
