package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/entities"
	"gorm.io/gorm"
)

// DBServiceError represents a generic Database Service error.
type DBServiceError struct {
	Msg string
	Err error
}

func (e *DBServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
	}
	return e.Msg
}
func (e *DBServiceError) Unwrap() error {
	return e.Err
}

// DBDUPError represents a duplicate error.
type DBDUPError struct{}

func (e *DBDUPError) Error() string { return "database error: duplicate entry" }

// DBNotFoundError represents a not found operation error.
type DBNotFoundError struct{}

func (e *DBNotFoundError) Error() string { return "database error: entry not found" }

// DatabaseService represents the database service.
type DatabaseService struct {
	Database *Database
}

// NewDatabaseService returns a new DatabaseService.
func NewDatabaseService(host string, port int, username string, password string, dbname string) (dbs *DatabaseService, err error) {
	dbs = &DatabaseService{}
	dbs.Database, err = NewDatabase(host, port, username, password, dbname)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

// Close closes all database connections.
func (dbs *DatabaseService) Close() error {
	return dbs.Database.Close()
}

// HealthCheck checks whether the database is still around.
func (dbs *DatabaseService) HealthCheck() error {
	return dbs.Database.HealthCheck()
}

// GetArticles returns all articles records matching a certain criteria.
func (dbs *DatabaseService) GetArticles(provider string, category string, sorting string, limit int, after *time.Time) (articles entities.Articles, err error) {
	articleRecords, err := dbs.Database.FindAllArticleRecords(provider, category, sorting, limit, after)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Articles{}, nil
	} else if err != nil {
		return nil, &DBServiceError{Msg: "database error", Err: err}
	}

	articleList := make(entities.Articles, 0, len(articleRecords))

	for _, articleRecord := range articleRecords {
		articleItem := entities.Article{
			GUID:          articleRecord.GUID,
			Title:         articleRecord.Title,
			Description:   articleRecord.Description,
			Link:          articleRecord.Link,
			PublishedTime: articleRecord.PublishedDate,
			Provider:      articleRecord.Provider.Name,
			Category:      articleRecord.Category.Name,
		}

		articleList = append(articleList, articleItem)
	}

	return articleList, nil
}

// AddArticle adds a new article record to the database.
func (dbs *DatabaseService) AddArticle(article entities.Article) (err error) {
	err = dbs.Database.InsertArticleRecord(article)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
				return &DBDUPError{}
			}
		}
		return &DBServiceError{Msg: "database error", Err: err}
	}

	return nil
}
