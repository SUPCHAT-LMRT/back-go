package search

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SearchTermHandler struct {
	useCase *SearchTermUseCase
}

func NewSearchTermHandler(useCase *SearchTermUseCase) *SearchTermHandler {
	return &SearchTermHandler{useCase: useCase}
}

func (h SearchTermHandler) Handle(c *gin.Context) {
	query := c.Query("q")
	kind := c.Query("kind")

	results, err := h.useCase.Execute(c, query, kind)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
