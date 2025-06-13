package list_mentionnable_user

import "github.com/gin-gonic/gin"

type ListMentionnableUserHandler struct {
	useCase *ListMentionnableUserUseCase
}

func NewListMentionnableUserHandler(useCase *ListMentionnableUserUseCase) *ListMentionnableUserHandler {
	return &ListMentionnableUserHandler{useCase: useCase}
}

func (h *ListMentionnableUserHandler) Handler(c *gin.Context) {
	channelId := c.Param("channelId")
	if channelId == "" {
		c.JSON(400, gin.H{"error": "channelId is required"})
		return
	}

	users, err := h.useCase.Execute(c.Request.Context(), channelId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"users": users})
}
