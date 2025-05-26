package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
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

	authenticatedUser := c.MustGet("user").(*user_entity.User)

	results, err := h.useCase.Execute(c, query, kind, authenticatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if results == nil {
		results = make([]*SearchResult, 0)
	}

	c.JSON(http.StatusOK, results)
}
