package list_workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
)

type DiscoverListWorkspaceHandler struct {
	useCase            *DiscoverListWorkspacesUseCase
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

func NewDiscoverListWorkspaceHandler(
	useCase *DiscoverListWorkspacesUseCase,
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase,
) *DiscoverListWorkspaceHandler {
	return &DiscoverListWorkspaceHandler{useCase: useCase, getUserByIdUseCase: getUserByIdUseCase}
}

// Handle récupère la liste des espaces de travail publics découvrables
// @Summary Liste des espaces de travail publics
// @Description Récupère tous les espaces de travail publics disponibles pour rejoindre
// @Tags workspace,discovery
// @Accept json
// @Produce json
// @Success 200 {array} map[string]string "Liste des espaces de travail publics"
// @Failure 401 {object} map[string]string "Utilisateur non authentifié"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des espaces de travail"
// @Router /api/workspaces/discover [get]
// @Security ApiKeyAuth
func (h DiscoverListWorkspaceHandler) Handle(c *gin.Context) {
	workspaces, err := h.useCase.Execute(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	result := make([]gin.H, len(workspaces))
	for i, workspace := range workspaces {
		ownerUser, err := h.getUserByIdUseCase.Execute(c, entity.UserId(workspace.OwnerId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		result[i] = gin.H{
			"id":                 workspace.Id,
			"name":               workspace.Name,
			"topic":              workspace.Topic,
			"ownerName":          ownerUser.FullName(),
			"membersCount":       workspace.MembersCount,
			"onlineMembersCount": workspace.OnlineMembersCount,
		}
	}

	c.JSON(http.StatusOK, result)
}
