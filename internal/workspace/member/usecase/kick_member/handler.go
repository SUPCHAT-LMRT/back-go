package kick_member

import (
	"net/http"

	"github.com/gin-gonic/gin"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity3 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

type KickGroupMemberHandler struct {
	UseCase *KickMemberUseCase
}

func NewKickGroupMemberHandler(useCase *KickMemberUseCase) *KickGroupMemberHandler {
	return &KickGroupMemberHandler{UseCase: useCase}
}

// Handle expulse un membre d'un espace de travail
// @Summary Expulsion d'un membre
// @Description Retire un utilisateur de l'espace de travail (requiert des permissions administratives)
// @Tags workspace,member
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param user_id path string true "ID de l'utilisateur à expulser"
// @Success 204 {string} string "Membre expulsé avec succès"
// @Failure 400 {object} map[string]string "Erreur lors de l'expulsion du membre"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour expulser un membre"
// @Failure 404 {object} map[string]string "Membre ou espace de travail non trouvé"
// @Router /api/workspaces/{workspace_id}/members/{user_id} [delete]
// @Security ApiKeyAuth
func (h *KickGroupMemberHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	memberId := c.Param("user_id") // Renommez en `member_id` si nécessaire dans la route.

	err := h.UseCase.Execute(
		c.Request.Context(),
		entity2.WorkspaceId(workspaceId),
		entity3.WorkspaceMemberId(memberId),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
