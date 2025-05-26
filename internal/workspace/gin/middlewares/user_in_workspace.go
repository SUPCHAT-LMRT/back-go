package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
)

type UserInWorkspaceMiddleware struct {
	getWorkspaceMember *get_workpace_member.GetWorkspaceMemberUseCase
}

func NewUserInWorkspaceMiddleware(
	getWorkspaceMember *get_workpace_member.GetWorkspaceMemberUseCase,
) *UserInWorkspaceMiddleware {
	return &UserInWorkspaceMiddleware{
		getWorkspaceMember: getWorkspaceMember,
	}
}

// Execute must be called in a middleware chain after the user middleware that sets the user in the context
func (a *UserInWorkspaceMiddleware) Execute(c *gin.Context) {
	loggedInUserInter, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	loggedInUser := loggedInUserInter.(*user_entity.User) //nolint:revive

	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	workspaceMember, err := a.getWorkspaceMember.Execute(
		c,
		entity.WorkspaceId(workspaceId),
		loggedInUser.Id,
	)
	if err != nil {
		if errors.Is(err, repository.ErrWorkspaceMemberNotFound) {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": "User not found in workspace"},
			)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Set("workspace_member", workspaceMember)

	c.Next()
}
