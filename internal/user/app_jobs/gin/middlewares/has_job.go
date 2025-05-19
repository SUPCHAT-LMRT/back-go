package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
	permissions2 "github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/permissions"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type HasJobPermissionsMiddleware struct {
	CheckPermissionJobUseCase *permissions2.CheckPermissionJobUseCase
}

func NewHasJobPermissionsMiddleware(
	jobRepository repository.JobRepository,
) *HasJobPermissionsMiddleware {
	return &HasJobPermissionsMiddleware{
		CheckPermissionJobUseCase: permissions2.NewCheckPermissionJobUseCase(jobRepository),
	}
}

func (h *HasJobPermissionsMiddleware) Execute(requiredPermission uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*user_entity.User)

		hasPermission, err := h.CheckPermissionJobUseCase.Execute(c.Request.Context(), user.Id.String(), requiredPermission)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(403, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
