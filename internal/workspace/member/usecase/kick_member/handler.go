package kick_member

import (
	"github.com/gin-gonic/gin"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity3 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"net/http"
)

type KickMemberHandler struct {
	UseCase *KickMemberUseCase
}

func NewKickMemberHandler(useCase *KickMemberUseCase) *KickMemberHandler {
	return &KickMemberHandler{UseCase: useCase}
}

func (h *KickMemberHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	memberId := c.Param("user_id") // Renommez en `member_id` si nécessaire dans la route.

	err := h.UseCase.Execute(c.Request.Context(), entity2.WorkspaceId(workspaceId), entity3.WorkspaceMemberId(memberId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
