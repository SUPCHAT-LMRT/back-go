package export_all_user_data

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ExportAllUserDataHandler struct {
	UseCase *ExportAllUserDataUseCase
}

func NewExportAllUserDataHandler(
	useCase *ExportAllUserDataUseCase,
) *ExportAllUserDataHandler {
	return &ExportAllUserDataHandler{UseCase: useCase}
}

// Handle récupère la liste des conversations récentes d'un utilisateur
// @Summary Lister les conversations récentes
// @Description Récupère toutes les conversations récentes de l'utilisateur authentifié
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {array} models.RecentChatResponse "Liste des conversations récentes"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/user/data/export [get]
// @Security ApiKeyAuth
func (h *ExportAllUserDataHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*entity.User) //nolint:revive

	exportedData, err := h.UseCase.Execute(c.Request.Context(), user.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="user_data_export.json"`)
	c.Data(200, "application/json", exportedData)
}
