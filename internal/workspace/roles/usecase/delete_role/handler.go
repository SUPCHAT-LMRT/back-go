package delete_role

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteRoleHandler struct {
	useCase *DeleteRoleUseCase
}

func NewDeleteRoleHandler(useCase *DeleteRoleUseCase) *DeleteRoleHandler {
	return &DeleteRoleHandler{useCase: useCase}
}

func (h DeleteRoleHandler) Handle(c *gin.Context) {
	roleId := c.Param("role_id")
	if roleId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_id is required"})
		return
	}

	err := h.useCase.Execute(c, roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
