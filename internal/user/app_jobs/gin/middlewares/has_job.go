package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type HasJobPermissionsMiddleware struct {
	jobRepo repository.JobRepository
}

func NewHasJobPermissionsMiddleware(jobRepo repository.JobRepository) *HasJobPermissionsMiddleware {
	return &HasJobPermissionsMiddleware{jobRepo: jobRepo}
}

func (h *HasJobPermissionsMiddleware) Execute(requiredPermission uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobId := c.Param("job_id")
		if jobId == "" {
			c.JSON(400, gin.H{"error": "job_id is required"})
			c.Abort()
			return
		}

		job, err := h.jobRepo.FindById(c.Request.Context(), jobId)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if job == nil {
			c.JSON(404, gin.H{"error": "Job not found"})
			c.Abort()
			return
		}

		if !job.HasPermission(requiredPermission) {
			c.JSON(403, gin.H{
				"error":        "Forbidden",
				"displayError": "Vous n'avez pas la permission pour ce job.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
