package login

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
	"os"
)

type LoginHandler struct {
	loginUserUseCase *LoginUserUseCase
}

func NewLoginHandler(loginUserUseCase *LoginUserUseCase) *LoginHandler {
	return &LoginHandler{loginUserUseCase: loginUserUseCase}
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required,min=3"`
	RememberMe bool   `json:"rememberMe"`
}

type LoginUserResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
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

	response, err := l.loginUserUseCase.Execute(c, LoginUserRequest{
		Email:      request.Email,
		Password:   request.Password,
		RememberMe: request.RememberMe,
	})
	if err != nil {
		if errors.Is(err, InvalidUsernameOrPasswordErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          err.Error(),
				"messageDisplay": "Email ou mot de passe incorrect",
			})
			return
		}

		if errors.Is(err, UserNotVerifiedErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          err.Error(),
				"messageDisplay": "Votre compte n'est pas encore vérifié, veuillez vérifier votre boîte de réception",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "An error occurred while logging in",
		})
		return
	}

	c.SetCookie("accessToken", response.AccessToken, int(response.AccessTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)
	c.SetCookie("refreshToken", response.RefreshToken, int(response.RefreshTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)

	c.JSON(http.StatusOK, l.LoginUserResponse(response.User))
}

func (l LoginHandler) LoginUserResponse(user *entity.User) LoginUserResponse {
	return LoginUserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
