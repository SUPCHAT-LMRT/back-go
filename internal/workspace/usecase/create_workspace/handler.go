package create_workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

type CreateWorkspaceHandler struct {
	useCase *CreateWorkspaceUseCase
}

func NewCreateWorkspaceHandler(useCase *CreateWorkspaceUseCase) *CreateWorkspaceHandler {
	return &CreateWorkspaceHandler{useCase: useCase}
}

// Handle crée un nouvel espace de travail et y ajoute l'utilisateur comme propriétaire
// @Summary Création d'un espace de travail
// @Description Crée un nouvel espace de travail avec l'utilisateur connecté comme propriétaire
// @Tags workspace
// @Accept json
// @Produce json
// @Param body body object true "Informations de l'espace de travail"
// @Success 201 {object} map[string]string "Espace de travail créé avec succès"
// @Failure 400 {object} map[string]string "Requête invalide ou données incomplètes"
// @Failure 401 {object} map[string]string "Utilisateur non authentifié"
// @Failure 500 {object} map[string]string "Erreur lors de la création de l'espace de travail"
// @Router /api/workspaces [post]
// @Security ApiKeyAuth
func (l CreateWorkspaceHandler) Handle(c *gin.Context) {
	var body struct {
		Name  string `json:"name" binding:"required"`
		Topic string `json:"topic"`
		Type  string `json:"type" binding:"required,oneof=PUBLIC PRIVATE"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := userVal.(*user_entity.User) //nolint:revive

	workspace := entity.Workspace{
		Name:    body.Name,
		Topic:   body.Topic,
		Type:    l.mapWorkspaceType(body.Type),
		OwnerId: user.Id,
	}
	err := l.useCase.Execute(c, &workspace, &entity2.WorkspaceMember{
		UserId: user.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      workspace.Id,
		"name":    workspace.Name,
		"topic":   workspace.Topic,
		"type":    workspace.Type,
		"ownerId": workspace.OwnerId,
	})
}

func (l CreateWorkspaceHandler) mapWorkspaceType(typeStr string) entity.WorkspaceType {
	switch typeStr {
	case "PUBLIC":
		return entity.WorkspaceTypePublic
	case "PRIVATE":
		return entity.WorkspaceTypePrivate
	default:
		return entity.WorkspaceTypePrivate
	}
}
