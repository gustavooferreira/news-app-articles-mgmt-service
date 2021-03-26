package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gustavooferreira/news-app-articles-mgmt-service/pkg/core/log"
)

// TODO: We would replace this with a proper config library like Viper.

const AppPrefix = "NEWS_APP_ARTICLES_MGMT"

// Configuration holds the entire configuration
type Configuration struct {
	Webserver WebserverConfiguration
	Options   OptionsConfiguration
	Database  DatabaseConfiguration
}

// WebserverConfiguration holds configuration related to the webserver
type WebserverConfiguration struct {
	Host string
	Port int
}

// OptionsConfiguration holds general configuration
type OptionsConfiguration struct {
	// Development mode disables the panic recovery so we can see what was the actual problem.
	// and also, enables pprof
	DevMode bool

	LogLevel log.Level
}

// DatabaseConfiguration holds configuration related to the database
type DatabaseConfiguration struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

// NewConfig returns new default configuration
func NewConfig() (config Configuration) {
	config.setDefaults()
	return config
}

// LoadConfig loads and validates config (from env vars)
func (config *Configuration) LoadConfig() (err error) {

	if webserverHost, ok := os.LookupEnv(AppPrefix + "_WEBSERVER_HOST"); ok {
		config.Webserver.Host = webserverHost
	}

	if webserverPort, ok := os.LookupEnv(AppPrefix + "_WEBSERVER_PORT"); ok {
		config.Webserver.Port, err = strconv.Atoi(webserverPort)
		if err != nil || config.Webserver.Port <= 0 || config.Webserver.Port > 1<<16-1 {
			return fmt.Errorf("configuration error: [webserver port] input not allowed <%s>", webserverPort)
		}
	}

	if devMode, ok := os.LookupEnv(AppPrefix + "_OPTIONS_DEV_MODE"); ok {
		config.Options.DevMode, err = strconv.ParseBool(devMode)
		if err != nil {
			return fmt.Errorf("configuration error: [options devmode] unrecognizable boolean <%s>", devMode)
		}
	}

	if logLevel, ok := os.LookupEnv(AppPrefix + "_OPTIONS_LOG_LEVEL"); ok {
		config.Options.LogLevel, err = ParseLogLevel(logLevel)
		if err != nil {
			return fmt.Errorf("configuration error: [options loglevel] unrecognized log level")
		}
	}

	if dbHost, ok := os.LookupEnv(AppPrefix + "_DATABASE_HOST"); ok {
		config.Database.Host = dbHost
	} else {
		return fmt.Errorf("configuration error: [database host] mandatory config parameter missing")
	}

	if dbPort, ok := os.LookupEnv(AppPrefix + "_DATABASE_PORT"); ok {
		config.Database.Port, err = strconv.Atoi(dbPort)
		if err != nil || config.Database.Port <= 0 || config.Database.Port > 1<<16-1 {
			return fmt.Errorf("configuration error: [database port] input not allowed <%s>", dbPort)
		}
	}

	if dbUsername, ok := os.LookupEnv(AppPrefix + "_DATABASE_USERNAME"); ok {
		config.Database.Username = dbUsername
	} else {
		return fmt.Errorf("configuration error: [database username] mandatory config parameter missing")
	}

	if dbPassword, ok := os.LookupEnv(AppPrefix + "_DATABASE_PASSWORD"); ok {
		config.Database.Password = dbPassword
	} else {
		return fmt.Errorf("configuration error: [database password] mandatory config parameter missing")
	}

	if dbName, ok := os.LookupEnv(AppPrefix + "_DATABASE_DBNAME"); ok {
		config.Database.DBName = dbName
	} else {
		return fmt.Errorf("configuration error: [database dbname] mandatory config parameter missing")
	}

	return nil
}

// setDefaults sets the config default values.
func (config *Configuration) setDefaults() {
	// Webserver
	config.Webserver.Host = "127.0.0.1"
	config.Webserver.Port = 8080

	// Options
	config.Options.DevMode = false
	config.Options.LogLevel = log.INFO

	// Database
	config.Database.Port = 3306
}

// ParseLogLevel parses a string and returns a log level enum.
func ParseLogLevel(level string) (logLevel log.Level, err error) {
	level = strings.ToLower(level)

	switch level {
	case "debug":
		logLevel = log.DEBUG
	case "info":
		logLevel = log.INFO
	case "warning":
		logLevel = log.WARN
	case "error":
		logLevel = log.ERROR
	default:
		return 0, fmt.Errorf("log level unrecognised")
	}

	return logLevel, nil
}
