package create_attachment

import (
	"github.com/gin-gonic/gin"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	uberdig "go.uber.org/dig"
	"net/http"
)

type CreateChatDirectAttachmentHandlerDeps struct {
	uberdig.In
	CreateChannelMessageAttachmentUseCase *CreateChannelMessageAttachmentUseCase
}

type CreateChannelMessageAttachmentHandler struct {
	deps CreateChatDirectAttachmentHandlerDeps
}

func NewCreateChannelMessageAttachmentHandler(deps CreateChatDirectAttachmentHandlerDeps) *CreateChannelMessageAttachmentHandler {
	return &CreateChannelMessageAttachmentHandler{deps: deps}
}

func (h *CreateChannelMessageAttachmentHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "workspace_id is required"})
		return
	}

	channelId := c.Param("channel_id")
	if channelId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "channel_id is required"})
		return
	}

	workspaceMember := c.MustGet("workspace_member").(*workspace_member_entity.WorkspaceMember) //nolint:revive

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

	_, err = h.deps.CreateChannelMessageAttachmentUseCase.Execute(c, &CreateChannelMessageAttachmentInput{
		WorkspaceId:           entity.WorkspaceId(workspaceId),
		ChannelId:             channel_entity.ChannelId(channelId),
		SenderWorkspaceMember: workspaceMember,
		File:                  file,
		FileName:              fileHeader.Filename,
		ContentType:           fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
}
