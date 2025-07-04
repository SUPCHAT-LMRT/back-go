package list_recent_chats

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type ListRecentChatsHandlerDeps struct {
	uberdig.In
	UseCase        *ListRecentChatsUseCase
	ResponseMapper mapper.Mapper[*ListRecentChatsUseCaseOutput, *RecentChatResponse]
}

type ListRecentChatsHandler struct {
	deps ListRecentChatsHandlerDeps
}

func NewListRecentChatsHandler(deps ListRecentChatsHandlerDeps) *ListRecentChatsHandler {
	return &ListRecentChatsHandler{deps: deps}
}

// Handle récupère la liste des conversations récentes d'un utilisateur
// @Summary Lister les conversations récentes
// @Description Récupère toutes les conversations récentes de l'utilisateur authentifié
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {array} RecentChatResponse "Liste des conversations récentes"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/chats/recents [get]
// @Security ApiKeyAuth
func (h *ListRecentChatsHandler) Handle(c *gin.Context) {
	authenticatedUser := c.MustGet("user").(*user_entity.User) //nolint:revive

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
	Id          entity.RecentChatId            `json:"id"`
	Kind        entity.RecentChatKind          `json:"kind"`
	Name        string                         `json:"name"`
	LastMessage *RecentChatLastMessageResponse `json:"lastMessage"`
}

type RecentChatLastMessageResponse struct {
	Id         entity.RecentChatId `json:"id"`
	Content    string              `json:"content"`
	CreatedAt  time.Time           `json:"createdAt"`
	AuthorId   user_entity.UserId  `json:"authorId"`
	AuthorName string              `json:"authorName"`
}
