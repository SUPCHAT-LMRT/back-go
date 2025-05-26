package list_all_users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListUserHandler struct {
	UseCase ListUserUseCase
}

func NewListUserHandler(useCase ListUserUseCase) *ListUserHandler {
	return &ListUserHandler{UseCase: useCase}
}

func (h *ListUserHandler) Handle(c *gin.Context) {
	users, err := h.UseCase.Execute(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseUsers []ResponseUser
	for _, user := range users {
		responseUsers = append(responseUsers, ResponseUser{
			ID:        string(user.Id),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}

	c.JSON(http.StatusOK, responseUsers)
}

type ResponseUser struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
