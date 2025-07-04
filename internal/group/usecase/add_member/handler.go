package add_member

import (
	"net/http"

	"github.com/gin-gonic/gin"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	uberdig "go.uber.org/dig"
)

type AddMemberToGroupHandlerDeps struct {
	uberdig.In
	UseCase            *AddMemberToGroupUseCase
	GetUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

type AddMemberToGroupHandler struct {
	deps AddMemberToGroupHandlerDeps
}

func NewAddMemberToGroupHandler(deps AddMemberToGroupHandlerDeps) *AddMemberToGroupHandler {
	return &AddMemberToGroupHandler{deps: deps}
}

// Handle ajoute un membre à un groupe de discussion
// @Summary Ajouter un membre au groupe
// @Description Ajoute un utilisateur spécifique à un groupe de discussion
// @Tags group
// @Accept json
// @Produce json
// @Param group_id path string true "ID du groupe"
// @Param request body AddMemberToGroupRequest true "Informations de l'utilisateur à ajouter"
// @Success 202 "Membre ajouté avec succès"
// @Failure 400 {object} map[string]string "Erreur de paramètre ou impossibilité de s'ajouter soi-même"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/groups/{group_id}/members [post]
// @Security ApiKeyAuth
func (h *AddMemberToGroupHandler) Handle(c *gin.Context) {
	groupId := c.Param("group_id")
	var req AddMemberToGroupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inviter := c.MustGet("user").(*user_entity.User) //nolint:revive

	if inviter.Id == req.InviteeUserId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't add yourself to the group"})
		return
	}

	invitee, err := h.deps.GetUserByIdUseCase.Execute(c, req.InviteeUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.deps.UseCase.Execute(c, group_entity.GroupId(groupId), invitee.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}

type AddMemberToGroupRequest struct {
	InviteeUserId user_entity.UserId `json:"inviteeUserId" binding:"required"`
}
