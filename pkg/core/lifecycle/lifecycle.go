package lifecycle

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/log"
)

// TerminateHandler terminates the application.
// This function waits on a SIGINT or SIGTERM signal and shuts down the HTTP server gracefully.
func TerminateHandler(logger log.Logger, server core.ShutDowner) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down application ...")

	// We will wait 5 seconds for the server to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.ShutDown(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("server failed to shutdown gracefully: %s", err.Error()))
	}
}
