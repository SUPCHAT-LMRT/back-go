package validate

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ValidateResetPasswordHandler struct {
	validateResetPasswordUseCase *ValidateResetPasswordUseCase
}

func NewValidateResetPasswordHandler(
	validateResetPasswordUseCase *ValidateResetPasswordUseCase,
) *ValidateResetPasswordHandler {
	return &ValidateResetPasswordHandler{validateResetPasswordUseCase: validateResetPasswordUseCase}
}

// Handle valide et applique une demande de réinitialisation de mot de passe
// @Summary Validation de la réinitialisation de mot de passe
// @Description Valide le token de réinitialisation et change le mot de passe de l'utilisateur
// @Tags account
// @Accept json
// @Produce json
// @Param request body validate.ValidateResetPasswordRequest true "Données de réinitialisation du mot de passe"
// @Success 204 {string} string "Mot de passe réinitialisé avec succès"
// @Failure 400 {object} map[string]string "Paramètres invalides ou mots de passe non concordants"
// @Failure 500 {object} map[string]string "Erreur lors de la validation du token ou de la mise à jour du mot de passe"
// @Router /api/account/reset-password/validate [post]
func (h *ValidateResetPasswordHandler) Handle(c *gin.Context) {
	var request ValidateResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Empty body",
				"message": "The body of the request is empty.",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate the request body
	if request.NewPassword != request.NewPasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Passwords do not match",
			"message":        "The new password and its confirmation do not match.",
			"messageDisplay": "Le nouveau mot de passe et sa confirmation ne correspondent pas.",
		})
		return
	}

	err := h.validateResetPasswordUseCase.Execute(c, request.Token, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":          err.Error(),
			"message":        "An error occurred while validating the account.",
			"messageDisplay": "Une erreur s'est produite lors de la validation de la demande de changement de mot de passe.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

type ValidateResetPasswordRequest struct {
	Token                   uuid.UUID `json:"token"                   binding:"required,uuid"`
	NewPassword             string    `json:"newPassword"             binding:"required,min=8"`
	NewPasswordConfirmation string    `json:"newPasswordConfirmation" binding:"required,eqfield=NewPassword"`
}
