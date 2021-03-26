package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// NoRoute provides a generic handler for unmatched routes.
func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "no route found"})
}

// RespondWithError is a helper function to return an error according to the API specification.
func RespondWithError(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, gin.H{"message": message})
}

// Healthcheck checks health of the service.
func (s *Server) Healthcheck(c *gin.Context) {
	err := s.Repo.HealthCheck()
	if err != nil {
		s.Logger.Error(fmt.Sprintf("database health check error: %s", err.Error()))
		c.JSON(500, gin.H{"status": "FAIL"})
		return
	}

	c.JSON(200, gin.H{"status": "OK"})
}
