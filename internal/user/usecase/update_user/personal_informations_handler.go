package update_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
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
		Email     string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request body",
		})
		return
	}

	userId := c.Param("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	user, err := h.useCase.GetUserById(c.Request.Context(), user_entity.UserId(userId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	user.FirstName = body.FirstName
	user.LastName = body.LastName
	user.Email = body.Email

	err = h.useCase.Execute(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":        "Failed to update user",
			"messageDisplay": "Une erreur est survenue lors de la mise à jour de vos informations personnelles. Veuillez réessayer plus tard.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
