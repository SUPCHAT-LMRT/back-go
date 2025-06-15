package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
)

type RequestResetPasswordHandler struct {
	userRepository              user_repository.UserRepository
	requestResetPasswordUseCase *RequestResetPasswordUseCase
}

func NewRequestResetPasswordHandler(
	userRepository user_repository.UserRepository,
	requestResetPasswordUseCase *RequestResetPasswordUseCase,
) *RequestResetPasswordHandler {
	return &RequestResetPasswordHandler{
		userRepository:              userRepository,
		requestResetPasswordUseCase: requestResetPasswordUseCase,
	}
}

// Handle envoie un email avec un lien de réinitialisation de mot de passe
// @Summary Demande de réinitialisation de mot de passe
// @Description Envoie un email avec un lien permettant à l'utilisateur de réinitialiser son mot de passe
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {string} string "Email envoyé avec succès"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur lors de l'envoi de l'email"
// @Router /api/account/reset-password/request [post]
// @Security ApiKeyAuth
func (h *RequestResetPasswordHandler) Handle(c *gin.Context) {
	loggedInUser, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":          "not_logged_in",
			"message":        "You are not logged in.",
			"messageDisplay": "Vous n'êtes pas connecté.",
		})
		return
	}

	user := loggedInUser.(*user_entity.User) //nolint:revive

	_, err := h.requestResetPasswordUseCase.Execute(c, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":          err.Error(),
			"message":        "An error occurred while sending the validation email.",
			"messageDisplay": "Une erreur s'est produite lors de l'envoi de l'email de validation.",
		})
		return
	}
}
