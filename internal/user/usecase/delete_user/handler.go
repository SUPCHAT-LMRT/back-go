package delete_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type DeleteUserHandler struct {
	useCase *DeleteUserUseCase
}

func NewDeleteUserHandler(useCase *DeleteUserUseCase) *DeleteUserHandler {
	return &DeleteUserHandler{useCase: useCase}
}

func (h *DeleteUserHandler) Handle(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	err := h.useCase.Execute(c.Request.Context(), user_entity.UserId(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
