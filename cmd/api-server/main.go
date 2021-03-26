package main

import (
	"fmt"
	"os"

	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/api"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/lifecycle"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/log"
	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/repository"
)

func main() {
	retCode := mainLogic()
	os.Exit(retCode)
}

func mainLogic() int {
	// Setup logger
	logger := core.NewAppLogger(os.Stdout, log.INFO)
	defer logger.Sync()

	logger.Info("APP starting")

	// Read config
	logger.Info("reading configuration", log.Field("type", "setup"))
	config := core.NewConfig()
	if err := config.LoadConfig(); err != nil {
		logger.Error(err.Error(), log.Field("type", "config"))
		return 1
	}

	// TODO: Set log level after reading config
	// something like this:
	// logger.SetLevel(config.Options.LogLevel)

	// Setup Database
	db, err := repository.NewDatabaseService(config.Database.Host, config.Database.Port,
		config.Database.Username, config.Database.Password, config.Database.DBName)
	if err != nil {
		logger.Error(fmt.Sprintf("database error: %s", err.Error()), log.Field("type", "setup"))
		return 1
	}
	defer db.Close()

	server := api.NewServer(config.Webserver.Host, config.Webserver.Port, config.Options.DevMode, logger, db)

	// Spawn SIGINT/SIGTERM listener
	go lifecycle.TerminateHandler(logger, server)

	logger.Info("listenning for incoming requests", log.Field("type", "runtime"))
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(fmt.Sprintf("unexpected error while serving HTTP: %s", err))
		return 1
	}

	logger.Info("APP gracefully terminated")
	return 0
}
