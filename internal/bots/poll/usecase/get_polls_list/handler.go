package get_polls_list

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetPollsListHandler struct {
	usecase *GetPollsListUseCase
}

func NewGetPollsListHandler(usecase *GetPollsListUseCase) *GetPollsListHandler {
	return &GetPollsListHandler{usecase: usecase}
}

// Handle récupère la liste des sondages d'un espace de travail
// @Summary Lister les sondages
// @Description Récupère tous les sondages disponibles dans un espace de travail
// @Tags poll
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Success 200 {array} PollResponse "Liste des sondages"
// @Failure 400 {object} map[string]string "ID de l'espace de travail manquant"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/workspaces/{workspace_id}/poll [get]
// @Security ApiKeyAuth
func (h *GetPollsListHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	user := c.MustGet("user").(*user_entity.User) //nolint:revive
	userId := string(user.Id)

	polls, err := h.usecase.Execute(c, workspace_entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]*PollResponse, len(polls))
	for i, poll := range polls {
		options := make([]OptionResponse, len(poll.Options))
		for j, opt := range poll.Options {
			isVoted := contains(opt.Voters, userId)
			options[j] = OptionResponse{
				Id:      opt.Id,
				Text:    opt.Text,
				Votes:   opt.Votes,
				Voters:  opt.Voters,
				IsVoted: isVoted,
			}
		}
		response[i] = &PollResponse{
			Id:        poll.Id,
			Question:  poll.Question,
			Options:   options,
			CreatedBy: poll.CreatedBy,
		}
	}

	c.JSON(http.StatusOK, response)
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

type PollResponse struct {
	Id        string           `json:"id"`
	Question  string           `json:"question"`
	Options   []OptionResponse `json:"options"`
	CreatedBy string           `json:"createdby"`
	ExpiresAt time.Time        `json:"expiresat"`
	CreatedAt time.Time        `json:"createdat"`
}

type OptionResponse struct {
	Id      string   `json:"id"`
	Text    string   `json:"text"`
	Votes   int      `json:"votes"`
	Voters  []string `json:"voters"`
	IsVoted bool     `json:"is_voted"`
}
