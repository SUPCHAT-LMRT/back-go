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
