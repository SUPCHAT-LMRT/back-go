package list_workpace_members

import (
	"context"
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_public_status"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	uberdig "go.uber.org/dig"
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
