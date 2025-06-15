package update_banner

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type UpdateWorkspaceBannerHandler struct {
	useCase *UpdateWorkspaceBannerUseCase
}

func NewUpdateWorkspaceBannerHandler(
	useCase *UpdateWorkspaceBannerUseCase,
) *UpdateWorkspaceBannerHandler {
	return &UpdateWorkspaceBannerHandler{useCase: useCase}
}

// Handle met à jour la bannière d'un espace de travail
// @Summary Mise à jour de la bannière
// @Description Modifie l'image de bannière d'un espace de travail
// @Tags workspace
// @Accept multipart/form-data
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param image formData file true "Fichier image pour la bannière"
// @Success 200 {object} map[string]string "Bannière mise à jour avec succès"
// @Failure 400 {object} map[string]string "Paramètres manquants ou invalides"
// @Failure 403 {object} map[string]string "Utilisateur non autorisé à gérer les paramètres de l'espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la mise à jour de la bannière"
// @Router /api/workspaces/{workspace_id}/banner [put]
// @Security ApiKeyAuth
func (l UpdateWorkspaceBannerHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "workspace_id is required",
		})
		return
	}

	// Get image from request (multipart/form-data)
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image is required",
		})
		return
	}

	// Open the file
	imageReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to open image",
		})
		return
	}

	err = l.useCase.Execute(c, entity.WorkspaceId(workspaceId), UpdateImage{
		ImageReader: imageReader,
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to update workspace banner",
		})
		return
	}

	c.Status(http.StatusOK)
}
