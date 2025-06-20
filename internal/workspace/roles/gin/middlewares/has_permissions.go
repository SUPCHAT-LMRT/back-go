package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/permissions"
)

type HasPermissionsMiddleware struct {
	checkPermissionUseCase *permissions.CheckPermissionUseCase
}

func NewHasPermissionsMiddleware(
	checkPermissionUseCase *permissions.CheckPermissionUseCase,
) *HasPermissionsMiddleware {
	return &HasPermissionsMiddleware{checkPermissionUseCase: checkPermissionUseCase}
}

func (h *HasPermissionsMiddleware) Execute(permissionsBit uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		workspaceMember := c.MustGet("workspace_member").(*workspace_entity.WorkspaceMember) //nolint:revive
		workspaceId := c.Param("workspace_id")
		if workspaceId == "" {
			c.JSON(400, gin.H{"error": "workspace_id is required"})
			c.Abort()
			return
		}

		hasPermission, err := h.checkPermissionUseCase.Execute(
			c,
			workspaceMember.Id,
			entity.WorkspaceId(workspaceId),
			permissionsBit,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "message": "Internal server error"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(403, gin.H{
				"error":        "Forbidden",
				"displayError": "Vous n'avez pas la permissions.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
