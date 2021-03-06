package core

import (
	"context"
	"time"

	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/entities"
)

// Repository represents a database holding the data
type Repository interface {
	HealthCheck() error
	GetArticles(provider string, category string, sorting string, limit int, after *time.Time) (articles entities.Articles, err error)
	AddArticle(article entities.Article) (err error)
}

// ShutDowner represents anything that can be shutdown like an HTTP server.
type ShutDowner interface {
	ShutDown(ctx context.Context) error
}
