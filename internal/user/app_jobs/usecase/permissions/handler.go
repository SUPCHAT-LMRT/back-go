package permissions

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckUserPermissionsHandler struct {
	checkPermissionUseCase *CheckPermissionJobUseCase
}

func NewCheckUserPermissionsHandler(
	useCase *CheckPermissionJobUseCase,
) *CheckUserPermissionsHandler {
	return &CheckUserPermissionsHandler{checkPermissionUseCase: useCase}
}

func (h *CheckUserPermissionsHandler) Handle(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	var request struct {
		Permissions uint64 `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hasPermission, err := h.checkPermissionUseCase.Execute(c, user_entity.UserId(userId), request.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hasPermission": hasPermission})
}
