package list_messages

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
	"net/http"
	"time"
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

func NewListDirectMessagesHandler(deps ListDirectMessagesHandlerDeps) *ListDirectMessagesHandler {
	return &ListDirectMessagesHandler{deps: deps}
}

func (h *ListDirectMessagesHandler) Handle(c *gin.Context) {
	authenticatedUser := c.MustGet("user").(*user_entity.User)

	otherUserId := c.Param("other_user_id")
	if otherUserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "other_user_id is required"})
		return
	}

	directMessages, err := h.deps.UseCase.Execute(c, authenticatedUser.Id, user_entity.UserId(otherUserId))
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

				reactionUsers[k] = DirectMessageReactionUserResponse{Id: userId.String(), Name: userReacted.FullName()}
			}

			reactions[j] = DirectMessageReactionResponse{
				Id:       reaction.Id.String(),
				Users:    reactionUsers,
				Reaction: reaction.Reaction,
			}
		}

		response[i] = DirectMessageResponse{
			Id:        message.Id.String(),
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
			Reactions: reactions,
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
	Id        string                          `json:"id"`
	Content   string                          `json:"content"`
	Author    DirectMessageAuthorResponse     `json:"author"`
	CreatedAt time.Time                       `json:"createdAt"`
	Reactions []DirectMessageReactionResponse `json:"reactions"`
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
