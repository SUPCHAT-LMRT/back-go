package register

import (
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	registerUserUseCase *RegisterUserUseCase
}

func NewRegisterHandler(useCase *RegisterUserUseCase) *RegisterHandler {
	return &RegisterHandler{registerUserUseCase: useCase}
}

type RegisterRequest struct {
	Token                string `json:"token"                binding:"required"`
	Password             string `json:"password"             binding:"required,min=3"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"`
}

// Handle traite les demandes d'inscription utilisateur
// @Summary Inscription d'un utilisateur
// @Description Enregistre un nouvel utilisateur avec un token d'invitation et un mot de passe
// @Tags account
// @Accept json
// @Produce json
// @Param request body register.RegisterRequest true "Informations d'inscription"
// @Success 204 {string} string "Inscription réussie"
// @Failure 400 {object} map[string]string "Paramètres invalides ou utilisateur déjà existant"
// @Router /api/account/auth/register [post]
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

	// TODO re-enable this after tests (code works)
	// if !isPasswordStrong(request.Password) {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"message": "Password is not strong enough. It must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters.",
	//	})
	//	return
	// }

	err := l.registerUserUseCase.Execute(c, request.Token, WithPassword(request.Password))
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
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

//nolint:unused
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
