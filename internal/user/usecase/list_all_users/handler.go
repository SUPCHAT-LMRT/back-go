package list_all_users

import (
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_public_status"
	uberdig "go.uber.org/dig"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListUserHandlerDeps struct {
	uberdig.In
	ListUserUseCase        *ListUserUseCase
	GetPublicStatusUseCase *get_public_status.GetPublicStatusUseCase
}

type ListUserHandler struct {
	deps ListUserHandlerDeps
}

func NewListUserHandler(deps ListUserHandlerDeps) *ListUserHandler {
	return &ListUserHandler{deps: deps}
}

func (h *ListUserHandler) Handle(c *gin.Context) {
	users, err := h.deps.ListUserUseCase.Execute(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseUsers []ResponseUser
	for _, user := range users {
		status, err := h.deps.GetPublicStatusUseCase.Execute(c, user.Id, entity.StatusOffline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to get user status"})
			return
		}

		responseUsers = append(responseUsers, ResponseUser{
			ID:        string(user.Id),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Status:    status,
		})
	}

	c.JSON(http.StatusOK, responseUsers)
}

type ResponseUser struct {
	ID        string        `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Email     string        `json:"email"`
	Status    entity.Status `json:"status"`
}
