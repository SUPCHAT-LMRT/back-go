package public_profile

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
	"net/http"
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

	c.JSON(http.StatusOK, PublicProfileResponse{
		Id:        profile.Id.String(),
		Email:     profile.Email,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Status:    profile.Status.String(),
	})
}

type PublicProfileResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Status    string `json:"status"`
}
