package register

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"regexp"
)

type RegisterHandler struct {
	registerUserUseCase *RegisterUserUseCase
}

func NewRegisterHandler(useCase *RegisterUserUseCase) *RegisterHandler {
	return &RegisterHandler{registerUserUseCase: useCase}
}

type RegisterRequest struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=3"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"`
}

func (l RegisterHandler) Handle(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Request body is empty.",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if request.Password != request.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password and password confirmation are not the same.",
		})
		return
	}

	if !isPasswordStrong(request.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password is not strong enough. It must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters.",
		})
		return
	}

	userRequest, err := l.RegisterUserRequest(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Cannot parse request",
		})
		return
	}

	err = l.registerUserUseCase.Execute(c, *userRequest)
	if err != nil {
		if errors.Is(err, UserAlreadyExistsErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          err.Error(),
				"messageDisplay": "Un utilisateur existe déjà avec cet email.",
				"level":          "warning",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func isPasswordStrong(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

	return hasLowercase && hasUppercase && hasDigit && hasSpecial
}

func (l RegisterHandler) RegisterUserRequest(request RegisterRequest) (*RegisterUserRequest, error) {
	return &RegisterUserRequest{
		Token:    request.Token,
		Password: request.Password,
	}, nil
}
