package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetArticles handles requests to get articles.
func (s *Server) GetArticles(c *gin.Context) {
	queryParams := struct {
		Provider string `form:"provider"`
		Category string `form:"category"`
	}{}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing query parameters: %s", err.Error()))
		RespondWithError(c, 400, err.Error())
		return
	}

	articles, err := s.Repo.GetArticles(queryParams.Provider, queryParams.Category)
	if err != nil {
		s.Logger.Error(err.Error())
		RespondWithError(c, 500, "Internal error")
		return
	}

	c.JSON(200, articles)
}

// AddArticle handles requests to add an article.
func (s *Server) AddArticle(c *gin.Context) {
}

// AddArticles handles requests to add multiple articles.
func (s *Server) AddArticles(c *gin.Context) {
}
