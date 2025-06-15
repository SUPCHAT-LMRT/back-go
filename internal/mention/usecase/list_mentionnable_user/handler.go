package list_mentionnable_user

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ListMentionnableUserHandler struct {
	useCase *ListMentionnableUserUseCase
}

func NewListMentionnableUserHandler(useCase *ListMentionnableUserUseCase) *ListMentionnableUserHandler {
	return &ListMentionnableUserHandler{useCase: useCase}
}

// Handle récupère la liste des utilisateurs mentionnables dans un canal
// @Summary Obtenir les utilisateurs mentionnables
// @Description Récupère la liste des utilisateurs qui peuvent être mentionnés dans un canal spécifique
// @Tags mention
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param channel_id path string true "ID du canal"
// @Success 200 {array} MentionnableUserResponse "Liste des utilisateurs mentionnables"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/workspaces/{workspace_id}/channels/{channel_id}/mentionnable-users [get]
// @Security ApiKeyAuth
func (h *ListMentionnableUserHandler) Handle(c *gin.Context) {
	channelId := c.Param("channel_id")
	if channelId == "" {
		c.JSON(400, gin.H{"error": "channelId is required"})
		return
	}

	users, err := h.useCase.Execute(c.Request.Context(), entity.ChannelId(channelId))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var mentionnableUsers []MentionnableUserResponse
	for _, user := range users {
		mentionnableUsers = append(mentionnableUsers, MentionnableUserResponse{
			Id:       user.Id.String(),
			Username: user.Username,
		})
	}

	c.JSON(200, mentionnableUsers)
}

type MentionnableUserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}
