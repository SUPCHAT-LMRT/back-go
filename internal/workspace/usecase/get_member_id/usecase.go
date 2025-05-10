package get_member_id

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetMemberIdHandler struct {
	Usecase GetMemberIdUsecase
}

func NewGetMemberIdHandler(usecase GetMemberIdUsecase) *GetMemberIdHandler {
	return &GetMemberIdHandler{Usecase: usecase}
}

func (h *GetMemberIdHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	userId := c.Param("user_id")

	if workspaceId == "" || userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id and user_id are required"})
		return
	}

	memberId, err := h.Usecase.Execute(c.Request.Context(), workspaceId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"member_id": memberId})
}
