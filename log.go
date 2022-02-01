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
	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel = go_logger.LOGGER_LEVEL_DEBUG
	case "INFO":
		logLevel = go_logger.LOGGER_LEVEL_INFO
	case "ERROR":
		logLevel = go_logger.LOGGER_LEVEL_ERROR
	case "EMERGENCY":
		logLevel = go_logger.LOGGER_LEVEL_EMERGENCY
	case "ALERT":
		logLevel = go_logger.LOGGER_LEVEL_ALERT
	case "CRITICAL":
		logLevel = go_logger.LOGGER_LEVEL_CRITICAL
	case "NOTICE":
		logLevel = go_logger.LOGGER_LEVEL_NOTICE
	case "WARNING":
		logLevel = go_logger.LOGGER_LEVEL_WARNING
	}

	consoleConfig := &go_logger.ConsoleConfig{
		Format: "%timestamp_format% [%level_string%] [%file%:%function%](line %line%): %body%",
	}
	GLogger.Attach("console", logLevel, consoleConfig)
}
