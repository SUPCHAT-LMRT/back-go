package get_my_account

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"net/http"
)

type GetMyUserAccountHandler struct {
	useCase *get_by_id.GetUserByIdUseCase
}

func NewGetMyUserAccountHandler(useCase *get_by_id.GetUserByIdUseCase) *GetMyUserAccountHandler {
	return &GetMyUserAccountHandler{useCase: useCase}
}

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Pseudo    string `json:"pseudo"`
}

func (g *GetMyUserAccountHandler) Handle(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, g.Response(user.(*entity.User)))
}

func (g *GetMyUserAccountHandler) Response(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:        user.Id.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Pseudo:    user.Pseudo,
	}
}
