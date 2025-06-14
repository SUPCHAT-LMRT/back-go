package export_all_user_data

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ExportAllUserDataHandler struct {
	UseCase *ExportAllUserDataUseCase
}

func NewExportAllUserDataHandler(
	useCase *ExportAllUserDataUseCase,
) *ExportAllUserDataHandler {
	return &ExportAllUserDataHandler{UseCase: useCase}
}

func (h *ExportAllUserDataHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*entity.User) //nolint:revive

	exportedData, err := h.UseCase.Execute(c.Request.Context(), user.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="user_data_export.json"`)
	c.Data(200, "application/json", exportedData)
}
