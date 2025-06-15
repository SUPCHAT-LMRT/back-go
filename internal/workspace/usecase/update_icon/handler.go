package update_icon

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type UpdateWorkspaceIconHandler struct {
	useCase *UpdateWorkspaceIconUseCase
}

func NewUpdateWorkspaceIconHandler(
	useCase *UpdateWorkspaceIconUseCase,
) *UpdateWorkspaceIconHandler {
	return &UpdateWorkspaceIconHandler{useCase: useCase}
}

// Handle met à jour l'icône d'un espace de travail
// @Summary Mise à jour de l'icône
// @Description Modifie l'image d'icône d'un espace de travail
// @Tags workspace
// @Accept multipart/form-data
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param image formData file true "Fichier image pour l'icône"
// @Success 200 {object} nil "Icône mise à jour avec succès"
// @Failure 400 {object} map[string]string "Paramètres manquants ou invalides"
// @Failure 403 {object} map[string]string "Utilisateur non autorisé dans cet espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la mise à jour de l'icône"
// @Router /api/workspaces/{workspace_id}/icon [put]
// @Security ApiKeyAuth
func (l UpdateWorkspaceIconHandler) Handle(c *gin.Context) {
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
			"message": "failed to update workspace icon",
		})
		return
	}

	c.Status(http.StatusOK)
}
