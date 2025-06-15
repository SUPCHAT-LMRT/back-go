package create_attachment

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
	"net/http"
)

type CreateChatDirectAttachmentHandlerDeps struct {
	uberdig.In
	CreateChatDirectAttachmentUseCase *CreateChatDirectAttachmentUseCase
}

type CreateChatDirectAttachmentHandler struct {
	deps CreateChatDirectAttachmentHandlerDeps
}

func NewCreateChatDirectAttachmentHandler(deps CreateChatDirectAttachmentHandlerDeps) *CreateChatDirectAttachmentHandler {
	return &CreateChatDirectAttachmentHandler{deps: deps}
}

func (h *CreateChatDirectAttachmentHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*user_entity.User) //nolint:revive

	otherUserId := c.Param("other_user_id")

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

	_, err = h.deps.CreateChatDirectAttachmentUseCase.Execute(c, &CreateChatDirectAttachmentInput{
		SenderUserId: user.Id,
		OtherUserId:  user_entity.UserId(otherUserId),
		File:         file,
		FileName:     fileHeader.Filename,
		ContentType:  fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
}
