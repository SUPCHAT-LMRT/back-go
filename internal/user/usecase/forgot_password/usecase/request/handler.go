package request

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models" // Import pour que Swagger trouve les modèles
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"net/http"
)

type RequestForgotPasswordHandler struct {
	userRepository               user_repository.UserRepository
	requestForgotPasswordUseCase *RequestForgotPasswordUseCase
}

func NewRequestForgotPasswordHandler(
	userRepository user_repository.UserRepository,
	requestForgotPasswordUseCase *RequestForgotPasswordUseCase,
) *RequestForgotPasswordHandler {
	return &RequestForgotPasswordHandler{
		userRepository:               userRepository,
		requestForgotPasswordUseCase: requestForgotPasswordUseCase,
	}
}

// Handle traite les demandes de récupération de mot de passe
// @Summary Demander une récupération de mot de passe
// @Description Envoie un email avec un lien pour réinitialiser le mot de passe d'un utilisateur
// @Tags account
// @Accept json
// @Produce json
// @Param request body models.RequestForgotPasswordRequest true "Email pour réinitialiser le mot de passe"
// @Success 200 {object} models.RequestForgotPasswordResponse "Email envoyé avec succès"
// @Failure 400 {object} map[string]string "Paramètres invalides"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/forgot-password/request [post]
func (h *RequestForgotPasswordHandler) Handle(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          err.Error(),
			"message":        "Please provide the required parameters",
			"messageDisplay": "Veuillez fournir les paramètres requis",
		})
		return
	}

	user, err := h.userRepository.GetByEmail(c, body.Email)
	if err != nil {
		return // if the user does not exist, we do not want to leak this information to the client
	}

	_, err = h.requestForgotPasswordUseCase.Execute(c, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":          err.Error(),
			"message":        "An error occurred while sending the validation email.",
			"messageDisplay": "Une erreur s'est produite lors de l'envoi de l'email de validation.",
		})
		return
	}
}
