package check_permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/permissions"
)

type CheckPermissionsHandler struct {
	checkPermissionUseCase *permissions.CheckPermissionUseCase
}

func NewCheckPermissionsHandler(
	checkPermissionUseCase *permissions.CheckPermissionUseCase,
) *CheckPermissionsHandler {
	return &CheckPermissionsHandler{checkPermissionUseCase: checkPermissionUseCase}
}

func (h *CheckPermissionsHandler) Handle(c *gin.Context) {
	workspaceMember := c.MustGet("workspace_member").(*workspace_entity.WorkspaceMember) //nolint:revive
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	var request struct {
		Permissions uint64 `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hasPermission, err := h.checkPermissionUseCase.Execute(
		c,
		workspaceMember.Id,
		entity2.WorkspaceId(workspaceId),
		request.Permissions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hasPermission": hasPermission})
}
