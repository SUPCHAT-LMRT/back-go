package list_all_users

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListUserHandler struct {
	UseCase ListUserUseCase
}

func NewListUserHandler(useCase ListUserUseCase) *ListUserHandler {
	return &ListUserHandler{UseCase: useCase}
}

func (h *ListUserHandler) Handle(c *gin.Context) {
	users, err := h.UseCase.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
