package get_minutely

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetMinutelyMessageSentHandler struct {
	useCase *GetMinutelyMessageSentUseCase
}

func NewGetMinutelyMessageSentHandler(
	useCase *GetMinutelyMessageSentUseCase,
) *GetMinutelyMessageSentHandler {
	return &GetMinutelyMessageSentHandler{useCase: useCase}
}

// Handle met à jour l'avatar d'un utilisateur
// @Summary Mise à jour de l'avatar utilisateur
// @Description Télécharge et associe une nouvelle image d'avatar à l'utilisateur connecté
// @Tags account
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Fichier image de l'avatar"
// @Success 200 {string} string "Avatar mis à jour avec succès"
// @Failure 400 {object} map[string]string "Paramètres invalides ou image manquante"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur lors du traitement de l'image ou de la mise à jour de l'avatar"
// @Router /api/account/avatar [patch]
// @Security ApiKeyAuth
func (h GetMinutelyMessageSentHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	messageSents, err := h.useCase.Execute(
		c,
		entity.WorkspaceId(workspaceId),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	messageSentsResponse := make([]*MessageSentResponse, len(messageSents))
	for i, sent := range messageSents {
		messageSentsResponse[i] = &MessageSentResponse{
			SentAt: sent.SentAt,
			Count:  sent.Count,
		}
	}

	c.JSON(http.StatusOK, messageSentsResponse)
}

type MessageSentResponse struct {
	SentAt time.Time `json:"sentAt"`
	Count  uint      `json:"count"`
}
