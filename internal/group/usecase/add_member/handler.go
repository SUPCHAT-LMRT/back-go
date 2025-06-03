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

func (h *AddMemberToGroupHandler) Handle(c *gin.Context) {
	var req AddMemberToGroupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inviter := c.MustGet("user").(*user_entity.User)

	if inviter.Id == req.InviteeUserId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't add yourself to the group"})
		return
	}

	invitee, err := h.deps.GetUserByIdUseCase.Execute(c, req.InviteeUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	group, err := h.deps.UseCase.Execute(c, req.GroupId, inviter.Id, invitee.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, group)
}

type AddMemberToGroupRequest struct {
	GroupId       *group_entity.GroupId `json:"groupId"`
	InviteeUserId user_entity.UserId    `json:"inviteeUserId" binding:"required"`
}
