package get_my_account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_or_create_status"
	uberdig "go.uber.org/dig"
)

type GetMyUserAccountHandlerDeps struct {
	uberdig.In
	GetOrCreateStatusUseCase *get_or_create_status.GetOrCreateStatusUseCase
}

type GetMyUserAccountHandler struct {
	deps GetMyUserAccountHandlerDeps
}

func NewGetMyUserAccountHandler(deps GetMyUserAccountHandlerDeps) *GetMyUserAccountHandler {
	return &GetMyUserAccountHandler{deps: deps}
}

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}

// Handle récupère les informations du compte de l'utilisateur authentifié
// @Summary Obtenir mon compte utilisateur
// @Description Récupère les informations du compte de l'utilisateur authentifié, y compris son statut actuel
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} get_my_account.UserResponse "Informations de l'utilisateur connecté"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/me [get]
// @Security ApiKeyAuth
func (g *GetMyUserAccountHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*entity.User) //nolint:revive

	userStatus, err := g.deps.GetOrCreateStatusUseCase.Execute(
		c,
		user.Id,
		user_status_entity.StatusOnline,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error(), "message": "Failed to save status"},
		)
		return
	}

	c.JSON(http.StatusOK, g.Response(user, userStatus))
}

func (g GetMyUserAccountHandler) Response(
	user *entity.User,
	status user_status_entity.Status,
) *UserResponse {
	return &UserResponse{
		ID:        user.Id.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Status:    status.String(),
	}
}
