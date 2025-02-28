package list_workpace_members

import (
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	uberdig "go.uber.org/dig"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type ListWorkspaceMembersHandlerDeps struct {
	uberdig.In
	UseCase            *ListWorkspaceMembersUseCase
	GetUserByIdUseCase *get_by_id.GetUserByIdUseCase
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

	totalMembers, members, err := h.deps.UseCase.Execute(c, entity.WorkspaceId(workspaceId), limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := make([]MemberResponse, len(members))
	for i, member := range members {
		username := member.Pseudo
		if username == "" {
			user, err := h.deps.GetUserByIdUseCase.Execute(c, member.UserId)
			if err != nil {
				continue
			}

			username = user.FullName()
		}

		result[i] = MemberResponse{
			Id:     string(member.Id),
			UserId: string(member.UserId),
			Pseudo: username,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"members": result,
		"total":   totalMembers,
	})
}

type MemberResponse struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Pseudo string `json:"pseudo"`
}
