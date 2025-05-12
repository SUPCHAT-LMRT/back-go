package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type HasJobPermissionsMiddleware struct {
	jobRepo repository.JobRepository
}

func NewHasJobPermissionsMiddleware(jobRepo repository.JobRepository) *HasJobPermissionsMiddleware {
	return &HasJobPermissionsMiddleware{jobRepo: jobRepo}
}

func (h *HasJobPermissionsMiddleware) Execute(requiredPermission uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*user_entity.User)

		job, err := h.jobRepo.FindById(c.Request.Context(), jobId)
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
		if job == nil {
			c.JSON(404, gin.H{"error": "job not found"})
			c.Abort()
			return
		}

		if !job.HasPermission(requiredPermission) {
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
