package update_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type UpdateAccountPersonalInformationsHandler struct {
	useCase *UpdateUserUseCase
}

func NewUpdateAccountPersonalInformationsHandler(
	useCase *UpdateUserUseCase,
) *UpdateAccountPersonalInformationsHandler {
	return &UpdateAccountPersonalInformationsHandler{useCase: useCase}
}

func (h *UpdateAccountPersonalInformationsHandler) Handle(c *gin.Context) {
	var body struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request body",
		})
		return
	}

	loggedInUser, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	user := loggedInUser.(*entity.User)

	user.FirstName = body.FirstName
	user.LastName = body.LastName

	err := h.useCase.Execute(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":          err.Error(),
			"message":        "Failed to update user",
			"messageDisplay": "Une erreur est survenue lors de la mise à jour de vos informations personnelles. Veuillez réessayer plus tard.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
