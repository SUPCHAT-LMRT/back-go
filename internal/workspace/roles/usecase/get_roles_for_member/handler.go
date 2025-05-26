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
