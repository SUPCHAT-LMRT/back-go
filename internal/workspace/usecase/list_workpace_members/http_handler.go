package list_workpace_members

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type ListWorkspaceMembersHandler struct {
	useCase *ListWorkspaceMembersUseCase
}

func NewListWorkspaceHandler(useCase *ListWorkspaceMembersUseCase) *ListWorkspaceMembersHandler {
	return &ListWorkspaceMembersHandler{useCase: useCase}
}

func (h *ListWorkspaceMembersHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspaceId")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "workspaceId is required",
		})
		return
	}

	members, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := make([]MemberResponse, len(members))
	for i, member := range members {
		result[i] = MemberResponse{
			Id:     string(member.Id),
			UserId: string(member.UserId),
			Pseudo: member.Pseudo,
		}
	}

	c.JSON(http.StatusOK, result)
}

type MemberResponse struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Pseudo string `json:"pseudo"`
}
