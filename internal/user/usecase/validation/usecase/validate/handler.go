package validate

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type ValidateAccountHandler struct {
	validateAccountUseCase *ValidateAccountUseCase
}

func NewValidateAccountHandler(validateAccountUseCase *ValidateAccountUseCase) *ValidateAccountHandler {
	return &ValidateAccountHandler{validateAccountUseCase: validateAccountUseCase}
}

func (h *ValidateAccountHandler) Handle(c *gin.Context) {
	var request ValidateAccountRequest
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

	err := h.validateAccountUseCase.Execute(c, request.ValidationToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "An error occurred while validating the account.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

type ValidateAccountRequest struct {
	ValidationToken uuid.UUID `json:"validationToken" binding:"required,uuid"`
}
