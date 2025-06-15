package public_profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type GetPublicProfileHandlerDeps struct {
	uberdig.In
	GetPublicUserProfileUseCase *GetPublicUserProfileUseCase
}

type GetPublicProfileHandler struct {
	deps GetPublicProfileHandlerDeps
}

func NewGetPublicProfileHandler(deps GetPublicProfileHandlerDeps) *GetPublicProfileHandler {
	return &GetPublicProfileHandler{deps: deps}
}

// Handle récupère le profil public d'un utilisateur
// @Summary Profil public d'un utilisateur
// @Description Récupère les informations publiques d'un utilisateur spécifique
// @Tags account
// @Accept json
// @Produce json
// @Param user_id path string true "ID de l'utilisateur"
// @Success 200 {object} public_profile.PublicProfileResponse "Informations du profil public"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/{user_id}/profile [get]
func (h *GetPublicProfileHandler) Handle(c *gin.Context) {
	userId := c.Param("user_id")

	profile, err := h.deps.GetPublicUserProfileUseCase.Execute(c, user_entity.UserId(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Unable to get user profile",
		})
		return
	}

	var jobsNames string
	for _, job := range profile.Jobs {
		if job.IsAssigned {
			if jobsNames != "" {
				jobsNames += ", "
			}
			jobsNames += job.Name
		}
	}

	c.JSON(http.StatusOK, PublicProfileResponse{
		Id:        profile.Id.String(),
		Email:     profile.Email,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Status:    profile.Status.String(),
		JobsNames: jobsNames,
	})
}

type PublicProfileResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Status    string `json:"status"`
	JobsNames string `json:"jobsNames"`
}
