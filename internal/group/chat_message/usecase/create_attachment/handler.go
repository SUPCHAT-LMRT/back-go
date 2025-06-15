package create_attachment

import (
	"github.com/gin-gonic/gin"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
	"net/http"
)

type CreateGroupAttachmentHandlerDeps struct {
	uberdig.In
	CreateGroupAttachmentUseCase *CreateGroupAttachmentUseCase
}

type CreateGroupAttachmentHandler struct {
	deps CreateGroupAttachmentHandlerDeps
}

func NewCreateGroupAttachmentHandler(deps CreateGroupAttachmentHandlerDeps) *CreateGroupAttachmentHandler {
	return &CreateGroupAttachmentHandler{deps: deps}
}

func (h *CreateGroupAttachmentHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*user_entity.User) //nolint:revive

	groupId := c.Param("group_id")

	// Récupération du fichier
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "No file received or file is invalid"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to open file"})
		return
	}
	defer file.Close()

	_, err = h.deps.CreateGroupAttachmentUseCase.Execute(c, &CreateGroupAttachmentInput{
		GroupId:      group_entity.GroupId(groupId),
		SenderUserId: user.Id,
		File:         file,
		FileName:     fileHeader.Filename,
		ContentType:  fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
}
