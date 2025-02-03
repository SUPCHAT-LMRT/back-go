package create_workspace

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"net/http"
)

type CreateWorkspaceHandler struct {
	useCase *CreateWorkspaceUseCase
}

func NewCreateWorkspaceHandler(useCase *CreateWorkspaceUseCase) *CreateWorkspaceHandler {
	return &CreateWorkspaceHandler{useCase: useCase}
}

func (l CreateWorkspaceHandler) Handle(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required,oneof=PUBLIC PRIVATE"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := userVal.(*user_entity.User)

	workspace := entity.Workspace{
		Name:    body.Name,
		Type:    l.mapWorkspaceType(body.Type),
		OwnerId: user.Id,
	}
	err := l.useCase.Execute(c, &workspace, &entity.WorkspaceMember{
		UserId: user.Id,
		Pseudo: user.Pseudo,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      workspace.Id,
		"name":    workspace.Name,
		"type":    workspace.Type,
		"ownerId": workspace.OwnerId,
	})
}

func (l CreateWorkspaceHandler) mapWorkspaceType(typeStr string) entity.WorkspaceType {
	switch typeStr {
	case "PUBLIC":
		return entity.WorkspaceTypePublic
	case "PRIVATE":
		return entity.WorkspaceTypePrivate
	default:
		return entity.WorkspaceTypePrivate
	}
}
