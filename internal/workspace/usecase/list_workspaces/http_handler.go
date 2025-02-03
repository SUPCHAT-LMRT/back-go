package list_workspaces

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type ListWorkspaceHandler struct {
	useCase *ListWorkspacesUseCase
}

func NewListWorkspaceHandler(useCase *ListWorkspacesUseCase) *ListWorkspaceHandler {
	return &ListWorkspaceHandler{useCase: useCase}
}

func (l ListWorkspaceHandler) Handle(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := userVal.(*user_entity.User)

	workspaces, err := l.useCase.Execute(c, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	result := make([]gin.H, len(workspaces))
	for i, workspace := range workspaces {
		result[i] = gin.H{
			"id":      workspace.Id,
			"name":    workspace.Name,
			"type":    workspace.Type,
			"ownerId": workspace.OwnerId,
		}
	}

	c.JSON(http.StatusOK, result)
}
