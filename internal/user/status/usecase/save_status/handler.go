package save_status

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models" // Import pour que Swagger trouve les modèles
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"net/http"
)

type SaveStatusHandler struct {
	useCase *SaveStatusUseCase
}

func NewSaveStatusHandler(useCase *SaveStatusUseCase) *SaveStatusHandler {
	return &SaveStatusHandler{useCase: useCase}
}

// Handle met à jour le statut d'un utilisateur
// @Summary Mettre à jour le statut d'un utilisateur
// @Description Met à jour le statut d'un utilisateur authentifié (en ligne, absent, ne pas déranger, etc.)
// @Tags account
// @Accept json
// @Produce json
// @Param request body models.SaveStatusRequest true "Informations du statut à mettre à jour"
// @Success 200 {object} models.SaveStatusResponse "Statut mis à jour avec succès"
// @Failure 400 {object} map[string]string "Paramètres invalides"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/status [patch]
// @Security ApiKeyAuth
func (h *SaveStatusHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*entity.User) //nolint:revive

	var body struct {
		Status string `json:"status" binding:"required,userStatus"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.Execute(c, user.Id, user_status_entity.Status(body.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}
