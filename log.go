package collycrawller

import (
	"os"
	"strings"

	go_logger "github.com/phachon/go-logger"
)

var (
	GLogger    *go_logger.Logger
	ModuleName = os.Getenv("MODULE_NAME")
	level      = os.Getenv("LOG_LEVEL")
)

func init() {
	GLogger = go_logger.NewLogger()

	// 日志等级默认为INFO
	logLevel := go_logger.LOGGER_LEVEL_INFO
	if strings.ToUpper(level) == "DEBUG" {
		logLevel = go_logger.LOGGER_LEVEL_DEBUG
	} else if strings.ToUpper(level) == "INFO" {
		logLevel = go_logger.LOGGER_LEVEL_INFO
	} else if strings.ToUpper(level) == "ERROR" {
		logLevel = go_logger.LOGGER_LEVEL_ERROR
	} else if strings.ToUpper(level) == "EMERGENCY" {
		logLevel = go_logger.LOGGER_LEVEL_EMERGENCY
	} else if strings.ToUpper(level) == "ALERT" {
		logLevel = go_logger.LOGGER_LEVEL_ALERT
	} else if strings.ToUpper(level) == "CRITICAL" {
		logLevel = go_logger.LOGGER_LEVEL_CRITICAL
	} else if strings.ToUpper(level) == "NOTICE" {
		logLevel = go_logger.LOGGER_LEVEL_NOTICE
	} else if strings.ToUpper(level) == "WARNING" {
		logLevel = go_logger.LOGGER_LEVEL_WARNING
	}
	consoleConfig := &go_logger.ConsoleConfig{
		Format: "%timestamp_format% [%level_string%] [%file%:%function%](line %line%): %body%",
	}
	GLogger.Detach("console")
	GLogger.Attach("console", logLevel, consoleConfig)
}
