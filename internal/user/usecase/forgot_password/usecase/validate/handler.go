package validate

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type ValidateForgotPasswordHandler struct {
	validateForgotPasswordUseCase *ValidateForgotPasswordUseCase
}

func NewValidateForgotPasswordHandler(validateForgotPasswordUseCase *ValidateForgotPasswordUseCase) *ValidateForgotPasswordHandler {
	return &ValidateForgotPasswordHandler{validateForgotPasswordUseCase: validateForgotPasswordUseCase}
}

func (h *ValidateForgotPasswordHandler) Handle(c *gin.Context) {
	var request ValidateForgotPasswordRequest
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

	err := h.validateForgotPasswordUseCase.Execute(c, request.Token, request.NewPassword)
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

type ValidateForgotPasswordRequest struct {
	Token                   uuid.UUID `json:"token" binding:"required,uuid"`
	NewPassword             string    `json:"newPassword" binding:"required"`
	NewPasswordConfirmation string    `json:"newPasswordConfirmation" binding:"required,eqfield=NewPassword"`
}
