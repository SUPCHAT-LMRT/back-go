package request

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	"net/http"
)

type RequestAccountValidationHandler struct {
	userRepository                  repository.UserRepository
	requestAccountValidationUseCase *RequestAccountValidationUseCase
}

func NewRequestAccountValidationHandler(userRepository repository.UserRepository, requestAccountValidationUseCase *RequestAccountValidationUseCase) *RequestAccountValidationHandler {
	return &RequestAccountValidationHandler{
		userRepository:                  userRepository,
		requestAccountValidationUseCase: requestAccountValidationUseCase,
	}
}

func (h *RequestAccountValidationHandler) Handle(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "User already verified or does not exist",
			"message":        "User already verified or does not exist",
			"messageDisplay": "Utilisateur déjà vérifié ou inexistant",
		})
		return // if the user does not exist, we do not want to leak this information to the client
	}

	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "User already verified or does not exist",
			"message":        "User already verified or does not exist",
			"messageDisplay": "Utilisateur déjà vérifié ou inexistant",
		})
		return
	}

	_, err = h.requestAccountValidationUseCase.Execute(c, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":          err.Error(),
			"message":        "An error occurred while sending the validation email.",
			"messageDisplay": "Une erreur s'est produite lors de l'envoi de l'email de validation.",
		})
		return
	}

}
