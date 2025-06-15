package list_workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ListWorkspaceHandler struct {
	useCase *ListWorkspacesUseCase
}

func NewListWorkspaceHandler(useCase *ListWorkspacesUseCase) *ListWorkspaceHandler {
	return &ListWorkspaceHandler{useCase: useCase}
}

// Handle récupère la liste des espaces de travail de l'utilisateur connecté
// @Summary Espaces de travail de l'utilisateur
// @Description Récupère tous les espaces de travail auxquels l'utilisateur appartient
// @Tags workspace
// @Accept json
// @Produce json
// @Success 200 {array} map[string]string "Liste des espaces de travail de l'utilisateur"
// @Failure 401 {object} map[string]string "Utilisateur non authentifié"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des espaces de travail"
// @Router /api/workspaces [get]
// @Security ApiKeyAuth
func (l ListWorkspaceHandler) Handle(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := userVal.(*user_entity.User) //nolint:revive

	workspaces, err := l.useCase.Execute(c, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	result := make([]gin.H, len(workspaces))
	for i, workspace := range workspaces {
		result[i] = gin.H{
			"id":      workspace.Id,
			"name":    workspace.Name,
			"type":    workspace.Type,
			"ownerId": workspace.OwnerId,
		}
	}

	c.JSON(http.StatusOK, result)
}
