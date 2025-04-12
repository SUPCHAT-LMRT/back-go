package login

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_or_create_status"
	uberdig "go.uber.org/dig"
	"net/http"
	"os"
)

type LoginHandlerDeps struct {
	uberdig.In
	LoginUserUseCase         *LoginUserUseCase
	GetOrCreateStatusUseCase *get_or_create_status.GetOrCreateStatusUseCase
}

type LoginHandler struct {
	deps LoginHandlerDeps
}

func NewLoginHandler(deps LoginHandlerDeps) *LoginHandler {
	return &LoginHandler{deps: deps}
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required,min=3"`
	RememberMe bool   `json:"rememberMe"`
}

type LoginUserResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Status    string `json:"status"`
}

func (l LoginHandler) Handle(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Please provide the required parameters",
		})
		return
	}

	response, err := l.deps.LoginUserUseCase.Execute(c, LoginUserRequest{
		Email:      request.Email,
		Password:   request.Password,
		RememberMe: request.RememberMe,
	})
	if err != nil {
		if errors.Is(err, InvalidUsernameOrPasswordErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          err.Error(),
				"messageDisplay": "Email ou mot de passe incorrect",
				"level":          "error",
			})
			return
		}

		if errors.Is(err, UserNotVerifiedErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          err.Error(),
				"messageDisplay": "Votre compte n'est pas encore vérifié, veuillez vérifier votre boîte de réception",
				"level":          "warning",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "An error occurred while logging in",
		})
		return
	}

	userStatus, err := l.deps.GetOrCreateStatusUseCase.Execute(c, response.User.Id, user_status_entity.StatusOnline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to save status"})
		return
	}

	c.SetCookie("accessToken", response.AccessToken, int(response.AccessTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)
	c.SetCookie("refreshToken", response.RefreshToken, int(response.RefreshTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)

	c.JSON(http.StatusOK, l.LoginUserResponse(response.User, userStatus))
}

func (l LoginHandler) LoginUserResponse(user *entity.User, status user_status_entity.Status) LoginUserResponse {
	return LoginUserResponse{
		Id:        user.Id.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    status.String(),
	}
}
