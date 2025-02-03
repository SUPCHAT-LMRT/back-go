package update_icon

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"net/http"
)

type UpdateWorkspaceIconHandler struct {
	useCase *UpdateWorkspaceIconUseCase
}

func NewUpdateWorkspaceIconHandler(useCase *UpdateWorkspaceIconUseCase) *UpdateWorkspaceIconHandler {
	return &UpdateWorkspaceIconHandler{useCase: useCase}
}

func (l UpdateWorkspaceIconHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspaceId")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "workspaceId is required",
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
