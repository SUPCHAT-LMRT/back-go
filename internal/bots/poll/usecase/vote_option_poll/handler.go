package vote_option_poll

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type VoteOptionPollHandler struct {
	usecase *VoteOptionPollUseCase
}

func NewVoteOptionPollHandler(usecase *VoteOptionPollUseCase) *VoteOptionPollHandler {
	return &VoteOptionPollHandler{usecase: usecase}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *VoteOptionPollHandler) Handle(c *gin.Context) {
	pollId := c.Param("poll_id")
	optionId := c.Param("option_id")
	user := c.MustGet("user").(*user_entity.User)

	if pollId == "" || optionId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "MISSING_PARAMETERS",
			Message: "poll_id et option_id sont requis",
		})
		return
	}

	err := h.usecase.Execute(c, pollId, optionId, string(user.Id))
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
		"message": "Vote enregistré avec succès",
	})
}
