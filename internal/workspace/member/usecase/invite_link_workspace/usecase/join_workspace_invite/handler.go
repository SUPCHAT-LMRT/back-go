package join_workspace_invite

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type JoinWorkspaceInviteHandler struct {
	useCase *JoinWorkspaceInviteUseCase
}

func NewJoinWorkspaceInviteHandler(useCase *JoinWorkspaceInviteUseCase) *JoinWorkspaceInviteHandler {
	return &JoinWorkspaceInviteHandler{useCase: useCase}
}

func (h *JoinWorkspaceInviteHandler) Handle(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token cannot be empty"})
		return
	}

	user := c.MustGet("user").(*user_entity.User)

	err := h.useCase.Execute(c, token, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
