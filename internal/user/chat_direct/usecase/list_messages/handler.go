package list_messages

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
)

type ListDirectMessagesHandlerDeps struct {
	uberdig.In
	UseCase                   *ListDirectMessagesUseCase
	GetUserByIdUseCase        *get_by_id.GetUserByIdUseCase
	GetWorkspaceMemberUseCase *get_workpace_member.GetWorkspaceMemberUseCase
}

type ListDirectMessagesHandler struct {
	deps ListDirectMessagesHandlerDeps
}

type MessageQuery struct {
	Limit           int       `form:"limit,default=20,max=100"`
	Before          time.Time `form:"before"`
	After           time.Time `form:"after"`
	AroundMessageId string    `form:"aroundMessageId"`
}

func NewListDirectMessagesHandler(deps ListDirectMessagesHandlerDeps) *ListDirectMessagesHandler {
	return &ListDirectMessagesHandler{deps: deps}
}

//nolint:revive

// Handle récupère les messages directs entre l'utilisateur authentifié et un autre utilisateur
// @Summary Lister les messages directs
// @Description Récupère la liste des messages directs échangés entre l'utilisateur authentifié et un autre utilisateur spécifié
// @Tags chats
// @Accept json
// @Produce json
// @Param other_user_id path string true "ID de l'utilisateur avec qui la conversation est partagée"
// @Param limit query int false "Nombre maximum de messages à récupérer (défaut: 20, max: 100)" default(20)
// @Param before query string false "Récupérer les messages avant cette date (format ISO8601)"
// @Param after query string false "Récupérer les messages après cette date (format ISO8601)"
// @Param aroundMessageId query string false "Récupérer les messages autour de cet ID de message"
// @Success 200 {array} list_messages.DirectMessageResponse "Liste des messages directs"
// @Failure 400 {object} map[string]string "Paramètres invalides"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/chats/direct/{other_user_id}/messages [get]
// @Security ApiKeyAuth
func (h *ListDirectMessagesHandler) Handle(c *gin.Context) {
	authenticatedUser := c.MustGet("user").(*user_entity.User)

	otherUserId := c.Param("other_user_id")
	if otherUserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "other_user_id is required"})
		return
	}

	var query MessageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	directMessages, err := h.deps.UseCase.Execute(
		c,
		authenticatedUser.Id,
		user_entity.UserId(otherUserId),
		QueryParams{
			Limit:           query.Limit,
			Before:          query.Before,
			After:           query.After,
			AroundMessageId: chat_direct_entity.ChatDirectId(query.AroundMessageId),
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := make([]DirectMessageResponse, len(directMessages))
	for i, message := range directMessages {
		reactions := make([]DirectMessageReactionResponse, len(message.Reactions))
		for j, reaction := range message.Reactions {
			reactionUsers := make([]DirectMessageReactionUserResponse, len(reaction.UserIds))
			for k, userId := range reaction.UserIds {
				userReacted, err := h.deps.GetUserByIdUseCase.Execute(c, userId)
				if err != nil {
					continue
				}

				reactionUsers[k] = DirectMessageReactionUserResponse{
					Id:   userId.String(),
					Name: userReacted.FullName(),
				}
			}

			reactions[j] = DirectMessageReactionResponse{
				Id:       reaction.Id.String(),
				Users:    reactionUsers,
				Reaction: reaction.Reaction,
			}
		}

		attachments := make([]DirectMessageAttachmentResponse, len(message.Attachments))
		for k, attachment := range message.Attachments {
			attachments[k] = DirectMessageAttachmentResponse{
				Id:   attachment.Id.String(),
				Name: attachment.FileName,
			}
		}

		response[i] = DirectMessageResponse{
			Id:          message.Id.String(),
			Content:     message.Content,
			CreatedAt:   message.CreatedAt,
			Reactions:   reactions,
			Attachments: attachments,
		}

		user, err := h.deps.GetUserByIdUseCase.Execute(c, message.SenderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "failed to get user by id",
			})
			return
		}

		response[i].Author = DirectMessageAuthorResponse{
			UserId:    user.Id.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}

	c.JSON(http.StatusOK, response)
}

type DirectMessageResponse struct {
	Id          string                            `json:"id"`
	Content     string                            `json:"content"`
	Author      DirectMessageAuthorResponse       `json:"author"`
	CreatedAt   time.Time                         `json:"createdAt"`
	Reactions   []DirectMessageReactionResponse   `json:"reactions"`
	Attachments []DirectMessageAttachmentResponse `json:"attachments"`
}

type DirectMessageAuthorResponse struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type DirectMessageReactionResponse struct {
	Id       string                              `json:"id"`
	Users    []DirectMessageReactionUserResponse `json:"users"`
	Reaction string                              `json:"reaction"`
}

type DirectMessageReactionUserResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DirectMessageAttachmentResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
