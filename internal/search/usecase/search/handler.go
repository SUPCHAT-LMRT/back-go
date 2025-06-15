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

// Handle recherche des éléments dans le système selon un terme
// @Summary Rechercher du contenu
// @Description Effectue une recherche globale basée sur un terme et filtrée par type de contenu
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Terme de recherche"
// @Param kind query string false "Type de contenu à rechercher (optionnel)"
// @Success 200 {array} SearchResult "Résultats de la recherche"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/search [get]
// @Security ApiKeyAuth
func (h SearchTermHandler) Handle(c *gin.Context) {
	query := c.Query("q")
	kind := c.Query("kind")

	authenticatedUser := c.MustGet("user").(*user_entity.User) //nolint:revive

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
