package api_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gustavooferreira/news-app-articles-mgmt-service/mocks"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/api"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/entities"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/log"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetArticlesHandler(t *testing.T) {
	assert := assert.New(t)

	logger := log.NullLogger{}
	mockDB := setupMockDB()
	server := api.NewServer("", 9999, false, logger, mockDB)
	router := server.Router

	baseURL := "/api/v1/articles"

	tests := map[string]struct {
		Provider           string
		Category           string
		Sorting            string
		Limit              int
		After              *time.Time
		expectedStatusCode int
	}{
		"test1": {
			Provider:           "errorCond",
			Category:           "errorCond",
			Sorting:            "desc",
			Limit:              50,
			expectedStatusCode: 500,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			rawURL := BuildQueryParams(baseURL, test.Provider, test.Category, test.Sorting, test.Limit, test.After)

			req, err := http.NewRequest("GET", rawURL, nil)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(test.expectedStatusCode, w.Code)
		})
	}

}

func BuildQueryParams(rawURL string, provider string, category string, sorting string, limit int, after *time.Time) string {
	v := url.Values{}

	if provider != "" {
		v.Set("provider", provider)
	}
	if category != "" {
		v.Set("category", category)
	}
	if sorting != "" {
		v.Set("sorting", sorting)
	}
	if limit != 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if after != nil {
		v.Set("after", after.String())
	}
	queryParams := v.Encode()

	if len(v) != 0 {
		rawURL += "?" + queryParams
	}

	return rawURL
}

func GenData() entities.Articles {
	data := entities.Articles{
		entities.Article{
			GUID:          "guid 1",
			Title:         "title 1",
			Description:   "descripion 1",
			Link:          "link 1",
			PublishedTime: time.Date(2020, 5, 10, 12, 30, 0, 0, time.UTC),
			Provider:      "provider 1",
			Category:      "category 1"},
		entities.Article{
			GUID:          "guid 2",
			Title:         "title 2",
			Description:   "descripion 2",
			Link:          "link 2",
			PublishedTime: time.Date(2020, 10, 10, 12, 30, 0, 0, time.UTC),
			Provider:      "provider 2",
			Category:      "category 2"},
		entities.Article{
			GUID:          "guid 3",
			Title:         "title 3",
			Description:   "descripion 3",
			Link:          "link 3",
			PublishedTime: time.Date(2020, 1, 12, 10, 0, 0, 0, time.UTC),
			Provider:      "provider 3",
			Category:      "category 3"},
	}

	return data
}

func setupMockDB() *mocks.Repository {
	mockDB := &mocks.Repository{}

	data := GenData()

	mockGetArticlesFn := func(provider string, category string, sorting string, limit int, after *time.Time) (articles entities.Articles) {
		articles = entities.Articles{}

		for _, item := range data {
			if provider != "" && provider != item.Provider {
				continue
			}

			if category != "" && category != item.Category {
				continue
			}

			articles = append(articles, item)
		}

		return articles
	}

	// GetArticles mock -------------------------------------
	// Error condition
	call := mockDB.On("GetArticles", "errorCond", "errorCond", "desc", 50, mock.Anything)
	call = call.Return(nil, &repository.DBServiceError{})

	// For every other case
	call = call.On("GetArticles", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	call = call.Return(mockGetArticlesFn, nil)

	return mockDB
}
