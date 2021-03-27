package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/entities"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/repository"
)

// GetArticles handles requests to get articles.
func (s *Server) GetArticles(c *gin.Context) {
	queryParams := struct {
		Provider string     `form:"provider"`
		Category string     `form:"category"`
		Sorting  string     `form:"sorting"`
		Limit    int        `form:"limit"`
		After    *time.Time `form:"after"`
	}{
		Sorting: "desc",
		Limit:   50,
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing query parameters: %s", err.Error()))
		RespondWithError(c, 400, err.Error())
		return
	}

	// Sorting can only be 'asc' or 'desc'
	if queryParams.Sorting != "asc" && queryParams.Sorting != "desc" {
		RespondWithError(c, 400, "sorting query parameter can only take one of two values: 'asc' or 'desc'")
		return
	}

	// Limit can have a max of 200
	if queryParams.Limit < 1 || queryParams.Limit > 200 {
		RespondWithError(c, 400, "limit query parameter can be a minimum of 1 and a maximum of 200")
		return
	}

	// Make timezone UTC
	if queryParams.After != nil {
		tempAfter := queryParams.After.UTC()
		queryParams.After = &tempAfter
	}

	articles, err := s.Repo.GetArticles(queryParams.Provider, queryParams.Category,
		queryParams.Sorting, queryParams.Limit, queryParams.After)
	if err != nil {
		s.Logger.Error(err.Error())
		RespondWithError(c, 500, "Internal error")
		return
	}

	c.JSON(200, articles)
}

// AddArticle handles requests to add an article.
func (s *Server) AddArticle(c *gin.Context) {
	bodyData := struct {
		GUID          string    `json:"guid" binding:"required"`
		Title         string    `json:"title" binding:"required"`
		Description   string    `json:"description" binding:"required"`
		Link          string    `json:"link" binding:"required"`
		PublishedTime time.Time `json:"published_date" binding:"required"`
		Provider      string    `json:"provider" binding:"required"`
		Category      string    `json:"category" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, err.Error())
		return
	}

	// Make timezone UTC
	bodyData.PublishedTime = bodyData.PublishedTime.UTC()

	article := entities.Article{
		GUID:          bodyData.GUID,
		Title:         bodyData.Title,
		Description:   bodyData.Description,
		Link:          bodyData.Link,
		PublishedTime: bodyData.PublishedTime,
		Provider:      bodyData.Provider,
		Category:      bodyData.Category,
	}

	err = s.Repo.AddArticle(article)
	if err, ok := err.(*repository.DBDUPError); ok {
		s.Logger.Error(err.Error())
		RespondWithError(c, 409, "article GUID already exists in the database")
		return
	} else if err != nil {
		s.Logger.Error(err.Error())
		RespondWithError(c, 500, "Internal error")
		return
	}

	c.Status(204)
}

// AddArticles handles requests to add multiple articles.
func (s *Server) AddArticles(c *gin.Context) {
	bodyData := []struct {
		GUID          string    `json:"guid" binding:"required"`
		Title         string    `json:"title" binding:"required"`
		Description   string    `json:"description" binding:"required"`
		Link          string    `json:"link" binding:"required"`
		PublishedTime time.Time `json:"published_date" binding:"required"`
		Provider      string    `json:"provider" binding:"required"`
		Category      string    `json:"category" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, err.Error())
		return
	}

	// convert into entities.Article
	for _, item := range bodyData {
		// Make timezone UTC
		item.PublishedTime = item.PublishedTime.UTC()

		article := entities.Article{
			GUID:          item.GUID,
			Title:         item.Title,
			Description:   item.Description,
			Link:          item.Link,
			PublishedTime: item.PublishedTime,
			Provider:      item.Provider,
			Category:      item.Category,
		}

		err = s.Repo.AddArticle(article)
		if err, ok := err.(*repository.DBDUPError); ok {
			continue
		} else if err != nil {
			s.Logger.Error(err.Error())
			RespondWithError(c, 500, "Internal error")
			return
		}
	}

	c.Status(204)
}
