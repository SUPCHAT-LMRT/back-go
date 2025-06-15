package create_group

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type CreateGroupHandler struct {
	useCase *CreateGroupUseCase
}

func NewCreateGroupHandler(useCase *CreateGroupUseCase) *CreateGroupHandler {
	return &CreateGroupHandler{
		useCase: useCase,
	}
}

// Handle crée un nouveau groupe de discussion
// @Summary Créer un groupe
// @Description Crée un nouveau groupe de discussion avec les utilisateurs spécifiés
// @Tags group
// @Accept json
// @Produce json
// @Param request body CreateGroupBody true "Informations du groupe à créer"
// @Success 200 {string} string "ID du groupe créé"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/groups [post]
// @Security ApiKeyAuth
func (h *CreateGroupHandler) Handle(c *gin.Context) {
	var input CreateGroupBody
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Invalid input"})
		return
	}

	user := c.MustGet("user").(*user_entity.User)

	// Convert UsersIds from string to UserId type if necessary
	usersIds := make([]user_entity.UserId, len(input.UsersIds))
	for i, id := range input.UsersIds {
		usersIds[i] = user_entity.UserId(id)
	}
	// Execute the use case
	group, err := h.useCase.Execute(c, CreateGroupInput{
		GroupName:   input.Name,
		UsersIds:    usersIds,
		OwnerUserId: user.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, group.Id)
}

type CreateGroupBody struct {
	Name     string   `json:"name" binding:"required"`
	UsersIds []string `json:"usersIds" binding:"required"`
}
