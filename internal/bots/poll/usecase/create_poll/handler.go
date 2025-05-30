package create_poll

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
	"time"
)

type CreatePollHandler struct {
	usecase *CreatePollUseCase
}

func NewCreatePollHandler(usecase *CreatePollUseCase) *CreatePollHandler {
	return &CreatePollHandler{usecase: usecase}
}

func (h *CreatePollHandler) Handle(c *gin.Context) {
	var req CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(*user_entity.User)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format"})
		return
	}

	poll, err := h.usecase.Execute(c, req.Question, req.Options, string(user.Id), workspaceId, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, poll)
}

type CreatePollRequest struct {
	Question  string   `json:"question" binding:"required"`
	Options   []string `json:"options" binding:"required,min=2"`
	ExpiresAt string   `json:"expiresat" binding:"required"`
}
