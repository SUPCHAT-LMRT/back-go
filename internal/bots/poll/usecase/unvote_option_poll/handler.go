package unvote_option_poll

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UnvoteOptionPollHandler struct {
	usecase *UnvoteOptionPollUseCase
}

func NewUnvoteOptionPollHandler(usecase *UnvoteOptionPollUseCase) *UnvoteOptionPollHandler {
	return &UnvoteOptionPollHandler{usecase: usecase}
}

func (h *UnvoteOptionPollHandler) Handle(c *gin.Context) {
	pollId := c.Param("poll_id")
	user := c.MustGet("user").(*user_entity.User) //nolint:revive

	if pollId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "MISSING_PARAMETERS",
			Message: "poll_id est requis",
		})
		return
	}

	err := h.usecase.Execute(c, pollId, string(user.Id))
	if err != nil {
		if customErr, ok := err.(*CustomError); ok {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    customErr.Code,
				Message: customErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote supprimé avec succès",
	})
}
