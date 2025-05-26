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
