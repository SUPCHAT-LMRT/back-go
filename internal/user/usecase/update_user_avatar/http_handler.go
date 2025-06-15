package update_user_avatar

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type UpdateUserAvatarHandler struct {
	useCase *UpdateUserAvatarUseCase
}

func NewUpdateUserAvatarHandler(useCase *UpdateUserAvatarUseCase) *UpdateUserAvatarHandler {
	return &UpdateUserAvatarHandler{useCase: useCase}
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
func (l UpdateUserAvatarHandler) Handle(c *gin.Context) {
	userInter, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user, ok := userInter.(*entity.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
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

	err = l.useCase.Execute(c, user.Id, UpdateAvatar{
		ImageReader: imageReader,
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to update user avatar",
		})
		return
	}

	c.Status(http.StatusOK)
}
