package leave_group

import (
	"errors"
	"github.com/gin-gonic/gin"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/kick_member"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type LeaveGroupHandler struct {
	useCase    *kick_member.KickMemberUseCase
	repository repository.GroupRepository
}

func NewLeaveGroupHandler(useCase *kick_member.KickMemberUseCase, repository repository.GroupRepository) *LeaveGroupHandler {
	return &LeaveGroupHandler{
		useCase:    useCase,
		repository: repository,
	}
}

func (h *LeaveGroupHandler) Handle(c *gin.Context) {
	groupId := group_entity.GroupId(c.Param("group_id"))

	user := c.MustGet("user").(*user_entity.User)

	groupMember, err := h.repository.GetMemberByUserId(c, groupId, user.Id)
	if err != nil {
		if errors.Is(err, repository.ErrGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.useCase.Execute(c, groupMember.Id, groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
