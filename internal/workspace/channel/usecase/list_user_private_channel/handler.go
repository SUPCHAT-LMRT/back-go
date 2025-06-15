package list_user_private_channel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_user_by_workspace_member_id"
	uberdig "go.uber.org/dig"
)

type ListPrivateChannelMembersHandlerDeps struct {
	uberdig.In
	ListPrivateChannelMembersUseCase  *ListPrivateChannelMembersUseCase
	GetUserByWorkspaceMemberIdUseCase *get_user_by_workspace_member_id.GetUserByWorkspaceMemberIdUseCase
}
type ListPrivateChannelMembersHandler struct {
	deps ListPrivateChannelMembersHandlerDeps
}

func NewListPrivateChannelMembersHandler(
	deps ListPrivateChannelMembersHandlerDeps,
) *ListPrivateChannelMembersHandler {
	return &ListPrivateChannelMembersHandler{deps: deps}
}

// Handle récupère la liste des membres d'un canal privé
// @Summary Liste des membres d'un canal privé
// @Description Retourne tous les utilisateurs qui ont accès à un canal privé spécifique
// @Tags workspace,channel,members
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param channel_id path string true "ID du canal privé"
// @Success 200 {array} Member "Liste des membres du canal privé"
// @Failure 400 {object} map[string]string "ID de canal manquant"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour accéder à ce canal"
// @Failure 404 {object} map[string]string "Canal non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des membres du canal"
// @Router /api/workspaces/{workspace_id}/channels/{channel_id}/members [get]
// @Security ApiKeyAuth
func (h *ListPrivateChannelMembersHandler) Handle(c *gin.Context) {
	channelId := c.Param("channel_id")
	membersIds, err := h.deps.ListPrivateChannelMembersUseCase.Execute(
		c.Request.Context(),
		entity.ChannelId(channelId),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var members []Member
	for _, memberId := range membersIds {
		user, err := h.deps.GetUserByWorkspaceMemberIdUseCase.Execute(c.Request.Context(), memberId)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to get user by ID: " + err.Error()},
			)
			return
		}
		members = append(members, Member{
			UserId:   user.Id.String(),
			Username: user.FullName(),
		})
	}
	c.JSON(http.StatusOK, members)
}

type Member struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}
