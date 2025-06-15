package unvote_option_poll

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UnvoteOptionPollHandler struct {
	usecase *UnvoteOptionPollUseCase
}

func NewUnvoteOptionPollHandler(usecase *UnvoteOptionPollUseCase) *UnvoteOptionPollHandler {
	return &UnvoteOptionPollHandler{usecase: usecase}
}

// Handle supprime un vote d'un utilisateur pour une option d'un sondage
// @Summary Supprimer un vote d'une option
// @Description Supprime le vote d'un utilisateur pour une option spécifique d'un sondage
// @Tags poll
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param poll_id path string true "ID du sondage"
// @Param option_id path string true "ID de l'option"
// @Success 200 {object} map[string]string "Vote supprimé avec succès"
// @Failure 400 {object} ErrorResponse "Erreur de paramètre ou vote non trouvé"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} ErrorResponse "Erreur interne du serveur"
// @Router /api/workspaces/{workspace_id}/poll/{poll_id}/unvote/{option_id} [post]
// @Security ApiKeyAuth
func (h *UnvoteOptionPollHandler) Handle(c *gin.Context) {
	pollId := c.Param("poll_id")
	user := c.MustGet("user").(*user_entity.User) //nolint:revive

	if pollId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "MISSING_PARAMETERS",
			Message: "poll_id est requis",
		})
		return
	}

	err := h.usecase.Execute(c, pollId, string(user.Id))
	if err != nil {
		if customErr, ok := err.(*CustomError); ok {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    customErr.Code,
				Message: customErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote supprimé avec succès",
	})
}
