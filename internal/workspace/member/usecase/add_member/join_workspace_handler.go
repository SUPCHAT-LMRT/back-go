package add_member

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace"
	uberdig "go.uber.org/dig"
	"net/http"
)

type AddMemberHandlerDeps struct {
	uberdig.In
	AddMemberUseCase    *AddMemberUseCase
	GetWorkspaceUseCase *get_workspace.GetWorkspaceUseCase
}

type AddMemberHandler struct {
	deps AddMemberHandlerDeps
}

func NewAddMemberHandler(deps AddMemberHandlerDeps) *AddMemberHandler {
	return &AddMemberHandler{deps: deps}
}

func (h *AddMemberHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*user_entity.User)
	workspaceId := c.Param("workspace_id")

	workspace, err := h.deps.GetWorkspaceUseCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if workspace.Type == entity.WorkspaceTypePrivate {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot join this workspace"})
		return
	}

	err = h.deps.AddMemberUseCase.Execute(c, entity.WorkspaceId(workspaceId), &workspace_entity.WorkspaceMember{
		WorkspaceId: entity.WorkspaceId(workspaceId),
		UserId:      user.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
