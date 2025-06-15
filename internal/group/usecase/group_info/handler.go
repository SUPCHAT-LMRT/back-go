package group_info

import (
	"github.com/gin-gonic/gin"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"net/http"
	"time"
)

type GetGroupInfoHandler struct {
	useCase *GetGroupInfoUseCase
}

func NewGetGroupInfoHandler(useCase *GetGroupInfoUseCase) *GetGroupInfoHandler {
	return &GetGroupInfoHandler{
		useCase: useCase,
	}
}

// Handle récupère les informations d'un groupe de discussion
// @Summary Obtenir les informations du groupe
// @Description Récupère les informations détaillées d'un groupe de discussion, y compris la liste de ses membres
// @Tags group
// @Accept json
// @Produce json
// @Param group_id path string true "ID du groupe"
// @Success 200 {object} GroupInfoResponse "Informations du groupe"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/groups/{group_id} [get]
// @Security ApiKeyAuth
func (h *GetGroupInfoHandler) Handle(c *gin.Context) {
	groupId := c.Param("group_id")

	info, err := h.useCase.Execute(c, group_entity.GroupId(groupId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to get group info"})
		return
	}

	response := &GroupInfoResponse{
		Id:        info.Id,
		Name:      info.Name,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
		Members:   make([]*GroupMemberResponse, len(info.Members)),
	}

	for i, member := range info.Members {
		response.Members[i] = &GroupMemberResponse{
			Id:           member.Id,
			UserId:       member.UserId,
			UserName:     member.UserName,
			Email:        member.Email,
			IsGroupOwner: member.IsGroupOwner,
			Status:       member.Status,
		}
	}

	c.JSON(http.StatusOK, response)
}

type GroupInfoResponse struct {
	Id        group_entity.GroupId   `json:"id"`
	Name      string                 `json:"name"`
	Members   []*GroupMemberResponse `json:"members"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type GroupMemberResponse struct {
	Id           group_entity.GroupMemberId `json:"id"`
	UserId       user_entity.UserId         `json:"userId"`
	UserName     string                     `json:"userName"`
	Email        string                     `json:"email"`
	IsGroupOwner bool                       `json:"isGroupOwner"`
	Status       entity.Status              `json:"status"`
}
