package list_recent_chats

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type ListRecentChatsHandlerDeps struct {
	uberdig.In
	UseCase        *ListRecentChatsUseCase
	ResponseMapper mapper.Mapper[*entity.RecentChat, *RecentChatResponse]
}

type ListRecentChatsHandler struct {
	deps ListRecentChatsHandlerDeps
}

func NewListRecentChatsHandler(deps ListRecentChatsHandlerDeps) *ListRecentChatsHandler {
	return &ListRecentChatsHandler{deps: deps}
}

func (h *ListRecentChatsHandler) Handle(c *gin.Context) {
	authenticatedUser := c.MustGet("user").(*user_entity.User)

	recentChats, err := h.deps.UseCase.Execute(c, authenticatedUser.Id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error(), "message": "failed to list recent chats"},
		)
		return
	}

	response := make([]*RecentChatResponse, len(recentChats))
	for i, recentChat := range recentChats {
		response[i], err = h.deps.ResponseMapper.MapToEntity(recentChat)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error(), "message": "failed to map recent chat"},
			)
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

type RecentChatResponse struct {
	Id   entity.RecentChatId   `json:"id"`
	Kind entity.RecentChatKind `json:"kind"`
	Name string                `json:"name"`
}
