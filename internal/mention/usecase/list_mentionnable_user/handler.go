package list_mentionnable_user

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ListMentionnableUserHandler struct {
	useCase *ListMentionnableUserUseCase
}

func NewListMentionnableUserHandler(useCase *ListMentionnableUserUseCase) *ListMentionnableUserHandler {
	return &ListMentionnableUserHandler{useCase: useCase}
}

func (h *ListMentionnableUserHandler) Handle(c *gin.Context) {
	channelId := c.Param("channel_id")
	if channelId == "" {
		c.JSON(400, gin.H{"error": "channelId is required"})
		return
	}

	users, err := h.useCase.Execute(c.Request.Context(), entity.ChannelId(channelId))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var mentionnableUsers []MentionnableUserResponse
	for _, user := range users {
		mentionnableUsers = append(mentionnableUsers, MentionnableUserResponse{
			Id:       user.Id.String(),
			Username: user.Username,
		})
	}

	c.JSON(200, mentionnableUsers)
}

type MentionnableUserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}
