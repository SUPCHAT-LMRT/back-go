package list_workspaces

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"net/http"
)

type DiscoverListWorkspaceHandler struct {
	useCase            *DiscoverListWorkspacesUseCase
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

func NewDiscoverListWorkspaceHandler(useCase *DiscoverListWorkspacesUseCase, getUserByIdUseCase *get_by_id.GetUserByIdUseCase) *DiscoverListWorkspaceHandler {
	return &DiscoverListWorkspaceHandler{useCase: useCase, getUserByIdUseCase: getUserByIdUseCase}
}

func (h DiscoverListWorkspaceHandler) Handle(c *gin.Context) {
	workspaces, err := h.useCase.Execute(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	result := make([]gin.H, len(workspaces))
	for i, workspace := range workspaces {
		ownerUser, err := h.getUserByIdUseCase.Execute(c, entity.UserId(workspace.OwnerId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		result[i] = gin.H{
			"id":                 workspace.Id,
			"name":               workspace.Name,
			"topic":              workspace.Topic,
			"ownerName":          ownerUser.FullName(),
			"membersCount":       workspace.MembersCount,
			"onlineMembersCount": workspace.OnlineMembersCount,
		}
	}

	c.JSON(http.StatusOK, result)
}
