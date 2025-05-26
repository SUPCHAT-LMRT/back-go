package dessassign_role

import "github.com/gin-gonic/gin"

type DessassignRoleFromUserHandler struct {
	useCase *DessassignRoleFromUserUsecase
}

func NewDessassignRoleFromUserHandler(
	useCase *DessassignRoleFromUserUsecase,
) *DessassignRoleFromUserHandler {
	return &DessassignRoleFromUserHandler{useCase: useCase}
}

func (h DessassignRoleFromUserHandler) Handle(c *gin.Context) {
	var req DessassignRoleFromUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.useCase.Execute(c, req.UserId, req.RoleId, req.WorkspaceId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Role unassigned successfully"})
}

type DessassignRoleFromUserRequest struct {
	RoleId      string `json:"role_id"      binding:"required,min=1,max=100"`
	UserId      string `json:"user_id"      binding:"required,min=1,max=100"`
	WorkspaceId string `json:"workspace_id" binding:"required,min=1,max=100"`
}
