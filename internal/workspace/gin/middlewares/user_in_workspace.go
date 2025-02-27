package middlewares

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/is_user_in_workspace"
	"net/http"
)

type UserInWorkspaceMiddleware struct {
	isUserInWorkspaceUseCase *is_user_in_workspace.IsUserInWorkspaceUseCase
}

func NewUserInWorkspaceMiddleware(isUserInWorkspaceUseCase *is_user_in_workspace.IsUserInWorkspaceUseCase) *UserInWorkspaceMiddleware {
	return &UserInWorkspaceMiddleware{isUserInWorkspaceUseCase: isUserInWorkspaceUseCase}
}

// Execute must be called in a middleware chain after the user middleware that sets the user in the context
func (a *UserInWorkspaceMiddleware) Execute(c *gin.Context) {
	loggedInUserInter, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	loggedInUser := loggedInUserInter.(*user_entity.User)

	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	isUserInWorkspace, err := a.isUserInWorkspaceUseCase.Execute(c, entity.WorkspaceId(workspaceId), loggedInUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isUserInWorkspace {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not in this workspace"})
		return
	}

	c.Next()
}
