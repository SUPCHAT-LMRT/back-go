package assign_role

import (
	"github.com/gin-gonic/gin"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
)

type AssignRoleToUserHandler struct {
	useCase *AssignRoleToUserUsecase
}

func NewAssignRoleToUserHandler(useCase *AssignRoleToUserUsecase) *AssignRoleToUserHandler {
	return &AssignRoleToUserHandler{useCase: useCase}
}
func (h AssignRoleToUserHandler) Handle(c *gin.Context) {
	var req AssignRoleToUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.useCase.Execute(
		c,
		entity2.WorkspaceMemberId(req.UserId),
		entity.RoleId(req.RoleId),
		entity.WorkspaceId(req.WorkspaceId),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Role assigned successfully"})
}

type AssignRoleToUserRequest struct {
	RoleId      string `json:"role_id" binding:"required,min=1,max=100"`
	UserId      string `json:"user_id" binding:"required,min=1,max=100"`
	WorkspaceId string `json:"workspace_id" binding:"required,min=1,max=100"`
}
