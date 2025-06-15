package get_member_id

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetMemberIdHandler struct {
	Usecase GetMemberIdUsecase
}

func NewGetMemberIdHandler(usecase GetMemberIdUsecase) *GetMemberIdHandler {
	return &GetMemberIdHandler{Usecase: usecase}
}

// Handle récupère l'ID d'un membre dans un espace de travail
// @Summary ID d'un membre
// @Description Obtient l'ID membre associé à un utilisateur dans un espace de travail spécifique
// @Tags workspace,member
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param user_id path string true "ID de l'utilisateur"
// @Success 200 {object} map[string]string "ID du membre trouvé"
// @Failure 400 {object} map[string]string "Paramètres manquants ou invalides"
// @Failure 403 {object} map[string]string "Utilisateur non autorisé dans cet espace de travail"
// @Failure 404 {object} map[string]string "Membre non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération de l'ID du membre"
// @Router /api/workspaces/{workspace_id}/members/{user_id} [get]
// @Security ApiKeyAuth
func (h *GetMemberIdHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	userId := c.Param("user_id")

	if workspaceId == "" || userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id and user_id are required"})
		return
	}

	memberId, err := h.Usecase.Execute(c.Request.Context(), workspaceId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"member_id": memberId})
}
