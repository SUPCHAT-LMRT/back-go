package join_workspace_invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type JoinWorkspaceInviteHandler struct {
	useCase *JoinWorkspaceInviteUseCase
}

func NewJoinWorkspaceInviteHandler(
	useCase *JoinWorkspaceInviteUseCase,
) *JoinWorkspaceInviteHandler {
	return &JoinWorkspaceInviteHandler{useCase: useCase}
}

// Handle permet à un utilisateur de rejoindre un espace de travail via un lien d'invitation
// @Summary Rejoindre via lien d'invitation
// @Description Permet à l'utilisateur authentifié de rejoindre un espace de travail en utilisant un lien d'invitation
// @Tags workspace,invite
// @Accept json
// @Produce json
// @Param token path string true "Token d'invitation"
// @Success 200 {string} string "Utilisateur ajouté à l'espace de travail"
// @Failure 400 {object} map[string]string "Token d'invitation manquant"
// @Failure 403 {object} map[string]string "Token d'invitation invalide ou expiré"
// @Failure 409 {object} map[string]string "L'utilisateur est déjà membre de cet espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de l'ajout de l'utilisateur à l'espace de travail"
// @Router /api/workspace-invite-link/{token}/join [post]
// @Security ApiKeyAuth
func (h *JoinWorkspaceInviteHandler) Handle(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token cannot be empty"})
		return
	}

	user := c.MustGet("user").(*user_entity.User) //nolint:revive

	err := h.useCase.Execute(c, token, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
