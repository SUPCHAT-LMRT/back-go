package list_workspace_members

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models" // Import pour que Swagger trouve les modèles
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_public_status"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	uberdig "go.uber.org/dig"
	"net/http"
	"strconv"
)

type ListWorkspaceMembersHandlerDeps struct {
	uberdig.In
	UseCase                *ListWorkspaceMembersUseCase
	GetUserByIdUseCase     *get_by_id.GetUserByIdUseCase
	GetPublicStatusUseCase *get_public_status.GetPublicStatusUseCase
}

type ListWorkspaceMembersHandler struct {
	deps ListWorkspaceMembersHandlerDeps
}

func NewListWorkspaceHandler(deps ListWorkspaceMembersHandlerDeps) *ListWorkspaceMembersHandler {
	return &ListWorkspaceMembersHandler{deps: deps}
}

// Handle récupère la liste des membres d'un espace de travail
// @Summary Liste des membres
// @Description Récupère la liste paginée des membres d'un espace de travail
// @Tags workspace,member
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param limit query int false "Nombre maximum de membres à retourner" default(10)
// @Param page query int false "Numéro de la page à retourner" default(1)
// @Success 200 {object} models.ListWorkspaceMembersResponse "Liste des membres de l'espace de travail"
// @Failure 400 {object} map[string]string "ID de workspace manquant ou paramètres de pagination invalides"
// @Failure 403 {object} map[string]string "Utilisateur non membre de l'espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des membres"
// @Router /api/workspaces/{workspace_id}/members [get]
// @Security ApiKeyAuth
func (h *ListWorkspaceMembersHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "workspace_id is required",
		})
		return
	}

	// Parse pagination parameters
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit: 10
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1")) // Default page: 1
	if err != nil || page < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	totalMembers, members, err := h.deps.UseCase.Execute(
		c,
		entity.WorkspaceId(workspaceId),
		limit,
		page,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"members": h.mapToMemberResponse(c, members),
		"total":   totalMembers,
	})
}

func (h *ListWorkspaceMembersHandler) mapToMemberResponse(
	ctx context.Context,
	members []*workspace_entity.WorkspaceMember,
) []MemberResponse {
	result := make([]MemberResponse, len(members))
	for i, member := range members {
		user, err := h.deps.GetUserByIdUseCase.Execute(ctx, member.UserId)
		if err != nil {
			continue
		}

		status, err := h.deps.GetPublicStatusUseCase.Execute(ctx, member.UserId, user_status_entity.StatusOffline)
		if err != nil {
			return nil
		}

		result[i] = MemberResponse{
			Id:     string(member.Id),
			UserId: string(member.UserId),
			Pseudo: user.FullName(),
			Email:  user.Email,
			Status: status.String(),
		}
	}

	return result
}

type MemberResponse struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Pseudo string `json:"pseudo"`
	Email  string `json:"email"`
	Status string `json:"status"`
}
