package list_messages

import (
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
)

type ListGroupMessagesHandlerDeps struct {
	uberdig.In
	UseCase                   *ListGroupChatMessagesUseCase
	GetUserByIdUseCase        *get_by_id.GetUserByIdUseCase
	GetWorkspaceMemberUseCase *get_workpace_member.GetWorkspaceMemberUseCase
}

type ListGroupChatMessagesHandler struct {
	deps ListGroupMessagesHandlerDeps
}

func NewListGroupChatMessagesHandler(
	deps ListGroupMessagesHandlerDeps,
) *ListGroupChatMessagesHandler {
	return &ListGroupChatMessagesHandler{deps: deps}
}

type MessageQuery struct {
	Limit           int       `form:"limit,default=20,max=100"`
	Before          time.Time `form:"before"`
	After           time.Time `form:"after"`
	AroundMessageId string    `form:"aroundMessageId"`
}

func (h *ListGroupChatMessagesHandler) Handle(c *gin.Context) {
	groupId := c.Param("group_id")
	if groupId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "group_id is required"})
		return
	}

	var query MessageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupChatMessages, err := h.deps.UseCase.Execute(c, group_entity.GroupId(groupId), QueryParams{
		Limit:           query.Limit,
		Before:          query.Before,
		After:           query.After,
		AroundMessageId: entity.GroupChatMessageId(query.AroundMessageId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := make([]GroupChatMessageResponse, len(groupChatMessages))
	for i, message := range groupChatMessages {
		reactions := make([]GroupMessageReactionResponse, len(message.Reactions))
		for j, reaction := range message.Reactions {
			reactionUsers := make([]GroupMessageReactionUserResponse, len(reaction.UserIds))
			for k, userId := range reaction.UserIds {
				userReacted, err := h.deps.GetUserByIdUseCase.Execute(c, userId)
				if err != nil {
					continue
				}

				reactionUsers[k] = GroupMessageReactionUserResponse{
					Id:   userId.String(),
					Name: userReacted.FullName(),
				}
			}

			reactions[j] = GroupMessageReactionResponse{
				Id:       reaction.Id.String(),
				Users:    reactionUsers,
				Reaction: reaction.Reaction,
			}
		}

		response[i] = GroupChatMessageResponse{
			Id:        message.Id.String(),
			GroupId:   message.GroupId.String(),
			Content:   message.Content,
			Reactions: reactions,
			CreatedAt: message.CreatedAt,
		}

		user, err := h.deps.GetUserByIdUseCase.Execute(c, message.AuthorId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "failed to get user by id",
			})
			return
		}

		response[i].Author = GroupMessageAuthorResponse{
			UserId:    user.Id.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}

	c.JSON(http.StatusOK, response)
}

type GroupChatMessageResponse struct {
	Id        string                         `json:"id"`
	GroupId   string                         `json:"groupId"`
	Content   string                         `json:"content"`
	Author    GroupMessageAuthorResponse     `json:"author"`
	Reactions []GroupMessageReactionResponse `json:"reactions"`
	CreatedAt time.Time                      `json:"createdAt"`
}

type GroupMessageAuthorResponse struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type GroupMessageReactionResponse struct {
	Id       string                             `json:"id"`
	Users    []GroupMessageReactionUserResponse `json:"users"`
	Reaction string                             `json:"reaction"`
}

type GroupMessageReactionUserResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
