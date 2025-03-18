package public_profile

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type GetPublicProfileHandler struct {
	useCase *GetPublicUserProfileUseCase
}

func NewGetPublicProfileHandler(useCase *GetPublicUserProfileUseCase) *GetPublicProfileHandler {
	return &GetPublicProfileHandler{useCase: useCase}
}

func (h *GetPublicProfileHandler) Handle(c *gin.Context) {
	userId := c.Param("user_id")

	profile, err := h.useCase.Execute(c, user_entity.UserId(userId))
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
	})
}

type PublicProfileResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
