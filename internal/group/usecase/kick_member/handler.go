package kick_member

import (
	"github.com/gin-gonic/gin"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	uberdig "go.uber.org/dig"
	"net/http"
)

type KickMemberHandlerDeps struct {
	uberdig.In
	KickMemberUseCase *KickMemberUseCase
}

type KickMemberHandler struct {
	deps KickMemberHandlerDeps
}

func NewKickMemberHandler(deps KickMemberHandlerDeps) *KickMemberHandler {
	return &KickMemberHandler{
		deps: deps,
	}
}

// Handle expulse un membre d'un groupe de discussion
// @Summary Expulser un membre du groupe
// @Description Expulse un utilisateur spécifique d'un groupe de discussion
// @Tags group
// @Accept json
// @Produce json
// @Param group_id path string true "ID du groupe"
// @Param request body models.KickMemberRequest true "ID du membre à expulser"
// @Success 200 "Membre expulsé avec succès"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/groups/{group_id}/members [delete]
// @Security ApiKeyAuth
func (h *KickMemberHandler) Handle(c *gin.Context) {
	groupId := group_entity.GroupId(c.Param("group_id"))
	var body struct {
		MemberId string `json:"memberId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.deps.KickMemberUseCase.Execute(c, group_entity.GroupMemberId(body.MemberId), groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to kick member"})
		return
	}

	c.Status(http.StatusOK)
}
